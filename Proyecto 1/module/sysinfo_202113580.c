#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/uaccess.h>
#include <linux/sched/signal.h>
#include <linux/seq_file.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Andres Agosto");
MODULE_DESCRIPTION("Modulo Kernel para Proyecto SO1");
MODULE_VERSION("1.0");

#define PROC_FILE "sysinfo_202113580"  // Cambiar al carnet correspondiente

static int sysinfo_show(struct seq_file *m, void *v) {
    // 1. Obtener información de memoria (a través de estructuras globales o funciones)
    // 2. Obtener información de CPU
    // 3. Iterar sobre tasks para filtrar contenedores Docker (podrías buscar "docker" en comm o cmdline)
    // 4. Construir una cadena con formato JSON

    seq_printf(m, "{\n");
    seq_printf(m, "\"mem_total\": %lu,\n", /* valor obtenido */);
    seq_printf(m, "\"mem_free\": %lu,\n",  /* valor obtenido */);
    seq_printf(m, "\"mem_used\": %lu,\n",  /* valor obtenido */);
    seq_printf(m, "\"cpu_usage\": %lu,\n", /* valor obtenido */);

    seq_printf(m, "\"containers\": [\n");
    // Por cada proceso "related to Docker"
    //   - PIDs, name, command line, % mem, % cpu, etc.
    seq_printf(m, "  {\n    \"pid\": %d,\n    \"name\": \"%s\", ... },\n");
    // ...
    seq_printf(m, "]\n");

    seq_printf(m, "}\n");

    return 0;
}

static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

static const struct file_operations fops = {
    .owner = THIS_MODULE,
    .open = sysinfo_open,
    .read = seq_read,
    .llseek = seq_lseek,
    .release = single_release,
};

static int __init sysinfo_init(void) {
    proc_create(PROC_FILE, 0, NULL, &fops);
    printk(KERN_INFO "sysinfo kernel module loaded\n");
    return 0;
}

static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_FILE, NULL);
    printk(KERN_INFO "sysinfo kernel module unloaded\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
