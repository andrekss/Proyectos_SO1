#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/string.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched.h>
#include <linux/timer.h>
#include <linux/jiffies.h>
#include <linux/uaccess.h>
#include <linux/tty.h>
#include <linux/sched/signal.h>
#include <linux/fs.h>
#include <linux/slab.h>
#include <linux/sched/mm.h>
#include <linux/binfmts.h>
#include <linux/timekeeping.h>
#include <linux/cgroup.h>
#include <linux/delay.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Andrés Agosto");
MODULE_DESCRIPTION("Modulo para leer y CPU en JSON");

#define ID_PREF "-id "
#define ID_MAX_LENGTH 100
#define BUFFER_SIZE 256

#define PROC_NAME "sysinfo_202113580"
#define CMDLine_Max_Lenght 256
void Liberar_Memoria(char *result);  // Prototipo de la función

// #define CONTAINER_ID_LENGTH 64
// Obtener Línea de comandos del proceso y retorno de apuntador del mismo
static char *get_process_cmdline(struct task_struct *task){

    struct mm_struct *mm;
    char *cmdline, *p;
    unsigned long arg_start, arg_end, env_start;
    int i, len;

    cmdline = kmalloc(CMDLine_Max_Lenght, GFP_KERNEL); // Reserva de memoria para la línea
    if (!cmdline)
        return NULL;

    mm = get_task_mm(task); // Info de memoria
    if (!mm)
    {
        kfree(cmdline);
        return NULL;
    }

    down_read(&mm->mmap_lock); // bloqueamos lectura mm_struct
    arg_start = mm->arg_start;
    arg_end = mm->arg_end;
    env_start = mm->env_start;
    up_read(&mm->mmap_lock);  // liberamos lectura mm_struct

    len = arg_end - arg_start;  // Longitud de la linea de comandos

    if (len > CMDLine_Max_Lenght - 1)
        len = CMDLine_Max_Lenght - 1;

    if (access_process_vm(task, arg_start, cmdline, len, 0) != len) // Los procesos solo pueden acceder a la memoria virtual
    {
        mmput(mm);
        kfree(cmdline);  // Obtenemos la memoria virtual de la linea de comandos del proceso
        return NULL;
    }
    cmdline[len] = '\0';

    p = cmdline;
    for (i = 0; i < len; i++)
        if (p[i] == '\0')
            p[i] = ' ';

    mmput(mm);
    return cmdline;
}

static char* extraer_id(const char *cmdline) {
    static char id[ID_MAX_LENGTH] = {0};
    char *pos = strstr(cmdline, ID_PREF);
    if (pos) {
        pos += strlen(ID_PREF);
        sscanf(pos, "%99s", id);
        return id;
    }
    return NULL;
}

// Procesos

void Liberar_Memoria(char *result) {
    if (result) {
        kfree(result);  // Liberar la memoria cuando ya no se necesite
    }
}

// Leer archivos en el kernel
static ssize_t Leer_Archivo(const char *path, char *buf, size_t max_size) {
    struct file *file;
    ssize_t bytes_read = 0;
    loff_t pos = 0; // Puntero del archivo

    file = filp_open(path, O_RDONLY, 0);
    if (IS_ERR(file)) return -ENOENT;

    // Leemos el archivo
    bytes_read = kernel_read(file, buf, max_size - 1, &pos);
    if (bytes_read < 0) {
        filp_close(file, NULL);
        return bytes_read;
    }

    buf[bytes_read] = '\0';
    filp_close(file, NULL);
    
    return bytes_read;
}

static char* get_Uso_Memoria(const char *container_id, unsigned long totalram) {
    char path[BUFFER_SIZE], buffer[BUFFER_SIZE];
    unsigned long long mem_usage = 0;
    unsigned long long porcent_mem_usage = 0;
    char *result;

    snprintf(path, sizeof(path), "/sys/fs/cgroup/system.slice/docker-%s.scope/memory.current", container_id);
    
    if (Leer_Archivo(path, buffer, sizeof(buffer)) > 0) {
        sscanf(buffer, "%llu", &mem_usage);
    }

    mem_usage /= (1024 * 1024);
    totalram /= 1024;

    if (totalram > 0) {
        porcent_mem_usage = (mem_usage * 10000) / totalram;
    }

    result = kmalloc(32, GFP_KERNEL);
    if (!result) return NULL;

    snprintf(result, 32, "%llu.%02llu", porcent_mem_usage / 100, porcent_mem_usage % 100);
    return result;
}

static char* get_Uso_CPU(const char *container_id) {
    char path[BUFFER_SIZE], buffer[BUFFER_SIZE];
    unsigned long long usage_usec_1 = 0, usage_usec_2 = 0;
    char *result;

    snprintf(path, sizeof(path), "/sys/fs/cgroup/system.slice/docker-%s.scope/cpu.stat", container_id);

    if (Leer_Archivo(path, buffer, sizeof(buffer)) <= 0) return NULL;
    sscanf(buffer, "usage_usec %llu", &usage_usec_1);

    msleep(700);  // Esperamos al cambio

    if (Leer_Archivo(path, buffer, sizeof(buffer)) <= 0) return NULL;
    sscanf(buffer, "usage_usec %llu", &usage_usec_2);

    unsigned long long delta_usage = usage_usec_2 - usage_usec_1;
    unsigned long long cpu_usage = (delta_usage) / 100;

    result = kmalloc(32, GFP_KERNEL);
    if (!result) return NULL;

    snprintf(result, 32, "%llu.%02llu", cpu_usage / 100, cpu_usage % 100);
    return result;
}

static char* get_Uso_Disco(const char *container_id) {
    char path[BUFFER_SIZE], buffer[BUFFER_SIZE];
    unsigned long long rbytes = 0, wbytes = 0, total_Uso = 0;
    char *result;

    snprintf(path, sizeof(path), "/sys/fs/cgroup/system.slice/docker-%s.scope/io.stat", container_id);
    
    if (Leer_Archivo(path, buffer, sizeof(buffer)) > 0) {
        char *rbytes_pos = strstr(buffer, "rbytes=");
        if (rbytes_pos) {
            rbytes_pos += strlen("rbytes=");
            sscanf(rbytes_pos, "%llu", &rbytes);
        }

        char *wbytes_pos = strstr(buffer, "wbytes=");
        if (wbytes_pos) {
            wbytes_pos += strlen("wbytes=");
            sscanf(wbytes_pos, "%llu", &wbytes);
        }
    }

    rbytes /= (1024 * 1024);
    wbytes /= (1024 * 1024);
    total_Uso = rbytes+wbytes;

    // Reservar memoria para el resultado
    result = kmalloc(64, GFP_KERNEL);
    if (!result) return NULL;

    snprintf(result, 64, "%llu", total_Uso); // MB

    return result;
}


static char* get_Uso_IO(const char *container_id) {
    char path[BUFFER_SIZE], buffer[BUFFER_SIZE];
    unsigned long long rios = 0, wios = 0, total_ops = 0;
    char *result;

    snprintf(path, sizeof(path), "/sys/fs/cgroup/system.slice/docker-%s.scope/io.stat", container_id);

    if (Leer_Archivo(path, buffer, sizeof(buffer)) > 0) {
        char *rios_pos = strstr(buffer, "rios=");
        if (rios_pos) {
            rios_pos += strlen("rios="); 
            sscanf(rios_pos, "%llu", &rios); 
        }

        char *wios_pos = strstr(buffer, "wios=");
        if (wios_pos) {
            wios_pos += strlen("wios=");
            sscanf(wios_pos, "%llu", &wios); 
        }
    }

    total_ops = rios+wios;

    result = kmalloc(64, GFP_KERNEL);
    if (!result) return NULL;

    snprintf(result, 64, "%llu", total_ops);

    return result;
}

static int sysinfo_show(struct seq_file *m, void *v) { // Mostrar en el proc
    struct sysinfo si;
    struct task_struct *task;  // recorrer procesos
    int first_process = 1;   // Saber primer proceso
    unsigned long total_jiffies = jiffies; // tiempo total cpu
    si_meminfo(&si);
    // Conversion a kB
    unsigned long totalram = si.totalram * (PAGE_SIZE / 1024); 
    unsigned long freeram = si.freeram * (PAGE_SIZE / 1024); 
    unsigned long ram_usada = totalram - freeram;

    seq_printf(m, "  {\n");
    seq_printf(m, "\"SystemInfo\": \n");
    seq_printf(m, "  {\n");
    seq_printf(m, "    \"RAM_Total\": %lu,\n", totalram);
    seq_printf(m, "    \"RAM_Libre\": %lu,\n", freeram);
    seq_printf(m, "    \"Uso_RAM\": %lu,\n", ram_usada);
    seq_printf(m, "    \"Uso_CPU\": %lu\n", total_jiffies);
    seq_printf(m, "  },\n");
    seq_printf(m, "\"Processes\": [\n");

    // Iteramos sobre los procesos
    for_each_process(task)
    {
        if (strcmp(task->comm, "containerd-shim") == 0) {
            unsigned long totalram = si.totalram * 4;
            char *cmdline = NULL;
            
            cmdline = get_process_cmdline(task);

            if (!first_process) {
                seq_printf(m, ",\n");
            } else {
                first_process = 0;
            }
            
            char *id_Contenedor = extraer_id(cmdline);
            char *disk_usage = get_Uso_Disco(id_Contenedor);
            char *cpu_usagep = get_Uso_CPU(id_Contenedor);
            char *mem_usagep = get_Uso_Memoria(id_Contenedor, totalram);
            char *ops_IO_usage = get_Uso_IO(id_Contenedor);

            seq_printf(m, "  {\n");
            seq_printf(m, "    \"PID\": %d,\n", task->pid);
            seq_printf(m, "    \"Name\": \"%s\",\n", task->comm);
            seq_printf(m, "    \"Cmdline\": \"%s\",\n", cmdline ? cmdline : "N/A");
            seq_printf(m, "    \"Porcentaje de uso Memoria\": %s,\n",mem_usagep );
            seq_printf(m, "    \"Porcentaje de uso CPU\": %s,\n", cpu_usagep);
            seq_printf(m, "    \"Uso de disco\": %s,\n", disk_usage);
            seq_printf(m, "    \"Operaciones I/O\": %s\n", ops_IO_usage);
            seq_printf(m, "  }");

            // liberamos la memoria
            Liberar_Memoria(disk_usage);
            Liberar_Memoria(mem_usagep);
            Liberar_Memoria(cpu_usagep);
            Liberar_Memoria(ops_IO_usage);
            // Liberamos la memoria de la línea de comandos
            if (cmdline) {
                kfree(cmdline);
            }
        }
    }

    seq_printf(m, "\n]\n}\n");
    return 0;
}

static int sysinfo_open(struct inode *inode, struct file *file)
{
    return single_open(file, sysinfo_show, NULL);
}

static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
};

static int __init sysinfo_init(void)
{
    proc_create(PROC_NAME, 0, NULL, &sysinfo_ops);
    printk(KERN_INFO "sysinfo_json modulo cargado\n");
    return 0;
}

static void __exit sysinfo_exit(void)
{
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "sysinfo_json modulo desinstalado\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);