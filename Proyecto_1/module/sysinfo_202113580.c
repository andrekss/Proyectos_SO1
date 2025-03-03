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

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Andrés Agosto");
MODULE_DESCRIPTION("Modulo para leer y CPU en JSON");

#define PROC_NAME "sysinfo_202113580"
#define CMDLine_Max_Lenght 256
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

static int sysinfo_show(struct seq_file *m, void *v) { // Mostrar en el proc
    struct sysinfo si;
    struct task_struct *task;  // recorrer procesos
    struct task_struct *hijos;
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
    seq_printf(m, "    \"Total_RAM\": %lu,\n", totalram);
    seq_printf(m, "    \"Free_RAM\": %lu,\n", freeram);
    seq_printf(m, "    \"Used_RAM\": %lu\n", ram_usada);
    seq_printf(m, "  },\n");
    seq_printf(m, "\"Processes\": [\n");

    // Iteramos sobre los procesos
    for_each_process(task)
    {
        if (strcmp(task->comm, "containerd-shim") == 0)
        {
            
            unsigned long vsz = 0;
            unsigned long rss = 0;
            //unsigned long totalram = si.totalram * 4;
            unsigned long mem_usage = 0;
            unsigned long cpu_usage = 0;
            char *cmdline = NULL;

            list_for_each_entry(hijos, &task->children, sibling){
                if (hijos->mm){
                    vsz = hijos->mm->total_vm << (PAGE_SHIFT - 10); // uso de vsz
                    rss = get_mm_rss(hijos->mm) << (PAGE_SHIFT - 10); // uso de rss
                    mem_usage = rss/totalram ;  // porcentaje de uso de memoria
                }

                unsigned long total_time = hijos->utime + hijos->stime + task->utime + task->stime;
                //cpu_usage = (total_time*100) / total_jiffies;
                cpu_usage = (total_time * 10000) / (total_jiffies );
                cmdline = get_process_cmdline(task);

                
                if (!first_process){
                    seq_printf(m, ",\n");
                }else{
                    first_process = 0;
                }

                seq_printf(m, "  {\n");
                seq_printf(m, "    \"PID\": %d,\n", task->pid);
                seq_printf(m, "    \"Name\": \"%s\",\n", task->comm);
                seq_printf(m, "    \"Cmdline\": \"%s\",\n", cmdline ? cmdline : "N/A");
                seq_printf(m, "    \"MemoryUsage\": %lu.%02lu,\n", mem_usage / 100, mem_usage % 100);
                seq_printf(m, "    \"CPUUsage\": %lu.%02lu\n", cpu_usage, cpu_usage % 100);
                seq_printf(m, "  }");

                // Liberamos la memoria de la línea de comandos
                if (cmdline){
                    kfree(cmdline);
                }
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