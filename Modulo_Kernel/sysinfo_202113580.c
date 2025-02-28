#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/uaccess.h>
#include <linux/sched/signal.h>
#include <linux/seq_file.h>
#include <linux/mm.h> 
#include <linux/sched.h> // recorre procesos

MODULE_LICENSE("GPL");
MODULE_AUTHOR("TuNombre");
MODULE_DESCRIPTION("Modulo Kernel para Proyecto SO1");
MODULE_VERSION("1.0");

#define PROC_FILE "sysinfo_202113580"

// Función para obtener la información del sistema y escribirla en /proc
static int sysinfo_show(struct seq_file *m, void *v) {
    struct sysinfo info;
    struct task_struct *task;
    
    si_meminfo(&info);
    unsigned long mem_total = info.totalram * 4 / 1024; // Convertir a MB
    unsigned long mem_free = info.freeram * 4 / 1024;
    unsigned long mem_used = mem_total - mem_free;

    seq_printf(m, "{\n");
    seq_printf(m, "\"mem_total\": %lu,\n", mem_total);
    seq_printf(m, "\"mem_free\": %lu,\n", mem_free);
    seq_printf(m, "\"mem_used\": %lu,\n", mem_used);
    seq_printf(m, "\"cpu_usage\": \"N/A\",\n"); // pendiente 

    seq_printf(m, "\"containers\": [\n");

    // Recorrer la lista de procesos y filtrar los contenedores Docker
    int first = 1;
    for_each_process(task) {
        if (strstr(task->comm, "docker") || strstr(task->comm, "containerd")) {
            if (!first) {
                seq_printf(m, ",\n");
            }
            first = 0;
            seq_printf(m, "  {\n    \"pid\": %d,\n    \"name\": \"%s\"\n  }", task->pid, task->comm);
        }
    }
    
    seq_printf(m, "\n]\n}\n"); // Cerrar Json

    return 0;
}

// Función que se ejecuta al abrir el archivo en /proc
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

static const struct proc_ops fops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

// Inicializar el modulo
static int __init sysinfo_init(void) {
    proc_create(PROC_FILE, 0, NULL, &fops);
    printk(KERN_INFO "sysinfo kernel module loaded\n");
    return 0;
}

// Limpieza del modulo
static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_FILE, NULL);
    printk(KERN_INFO "sysinfo kernel module unloaded\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
