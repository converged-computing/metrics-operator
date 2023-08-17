# LDMS Example

LDMS is typically run as a monitoring tool for nodes on an HPC cluster. You can
read [more about it here](https://github.com/ovis-hpc/ovis).

## Usage

Create a cluster and install JobSet to it.

```bash
kind create cluster
VERSION=v0.2.0
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
```

Install the operator (from the development manifest here):

```bash
$ kubectl apply -f ../../dist/metrics-operator-dev.yaml
```

How to see metrics operator logs:

```bash
$ kubectl logs -n metrics-system metrics-controller-manager-859c66464c-7rpbw
```

Then create the metrics set. This is going to run a simple sysstat tool to collect metrics
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be one):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS              RESTARTS   AGE
- metricset-sample-m-0-0-mkwrh   0/1     ContainerCreating   0          2m20s
+ metricset-sample-m-0-0-mkwrh   1/1     Running             0          3m10s
```

You can see logs for the metric, run twice (two completions) with 10 second breaks.

```bash
kubectl logs metricset-sample-m-0-czxrq 
```
```console
METADATA START {"pods":1,"completions":1,"metricName":"app-ldms","metricDescription":"provides LDMS, a low-overhead, low-latency framework for collecting, transferring, and storing metric data on a large distributed computer system.","metricType":"application","metricOptions":{"command":"ldms_ls -h localhost -x sock -p 10444 -l -v","completions":2,"rate":10,"workdir":"/opt"}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
Schema         Instance                 Flags  Msize  Dsize  Hsize  UID    GID    Perm       Update            Duration          Info    
-------------- ------------------------ ------ ------ ------ ------ ------ ------ ---------- ----------------- ----------------- --------
vmstat         metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/vmstat     L    8920   1384      0      0      0 -rwxr-xr-x          0.000000          0.000000 "updt_hint_us"="1000000:50000" 
meminfo        metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/meminfo     L    2888    520      0  12345  12345 -rwxr-xr-x          0.000000          0.000000 "updt_hint_us"="1000000:50000" 
-------------- ------------------------ ------ ------ ------ ------ ------ ------ ---------- ----------------- ----------------- --------
Total Sets: 2, Meta Data (kB): 11.81, Data (kB) 1.90, Memory (kB): 13.71

=======================================================================

metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/meminfo: inconsistent, last update: Thu Jan 01 00:00:00 1970 +0000 [0us] 
M u64          component_id                               1 
D u64          job_id                                     0 
D u64          app_id                                     0 
D u64          MemTotal                                   0 
D u64          MemFree                                    0 
D u64          MemAvailable                               0 
D u64          Buffers                                    0 
D u64          Cached                                     0 
D u64          SwapCached                                 0 
D u64          Active                                     0 
D u64          Inactive                                   0 
D u64          Active(anon)                               0 
D u64          Inactive(anon)                             0 
D u64          Active(file)                               0 
D u64          Inactive(file)                             0 
D u64          Unevictable                                0 
D u64          Mlocked                                    0 
D u64          SwapTotal                                  0 
D u64          SwapFree                                   0 
D u64          Dirty                                      0 
D u64          Writeback                                  0 
D u64          AnonPages                                  0 
D u64          Mapped                                     0 
D u64          Shmem                                      0 
D u64          KReclaimable                               0 
D u64          Slab                                       0 
D u64          SReclaimable                               0 
D u64          SUnreclaim                                 0 
D u64          KernelStack                                0 
D u64          PageTables                                 0 
D u64          NFS_Unstable                               0 
D u64          Bounce                                     0 
D u64          WritebackTmp                               0 
D u64          CommitLimit                                0 
D u64          Committed_AS                               0 
D u64          VmallocTotal                               0 
D u64          VmallocUsed                                0 
D u64          VmallocChunk                               0 
D u64          Percpu                                     0 
D u64          HardwareCorrupted                          0 
D u64          AnonHugePages                              0 
D u64          ShmemHugePages                             0 
D u64          ShmemPmdMapped                             0 
D u64          FileHugePages                              0 
D u64          FilePmdMapped                              0 
D u64          HugePages_Total                            0 
D u64          HugePages_Free                             0 
D u64          HugePages_Rsvd                             0 
D u64          HugePages_Surp                             0 
D u64          Hugepagesize                               0 
D u64          Hugetlb                                    0 
D u64          DirectMap4k                                0 
D u64          DirectMap2M                                0 
D u64          DirectMap1G                                0 

metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/vmstat: inconsistent, last update: Thu Jan 01 00:00:00 1970 +0000 [0us] 
M u64          component_id                               1 
D u64          job_id                                     0 
D u64          app_id                                     0 
D u64          nr_free_pages                              0 
D u64          nr_zone_inactive_anon                      0 
D u64          nr_zone_active_anon                        0 
D u64          nr_zone_inactive_file                      0 
D u64          nr_zone_active_file                        0 
D u64          nr_zone_unevictable                        0 
D u64          nr_zone_write_pending                      0 
D u64          nr_mlock                                   0 
D u64          nr_bounce                                  0 
D u64          nr_zspages                                 0 
D u64          nr_free_cma                                0 
D u64          numa_hit                                   0 
D u64          numa_miss                                  0 
D u64          numa_foreign                               0 
D u64          numa_interleave                            0 
D u64          numa_local                                 0 
D u64          numa_other                                 0 
D u64          nr_inactive_anon                           0 
D u64          nr_active_anon                             0 
D u64          nr_inactive_file                           0 
D u64          nr_active_file                             0 
D u64          nr_unevictable                             0 
D u64          nr_slab_reclaimable                        0 
D u64          nr_slab_unreclaimable                      0 
D u64          nr_isolated_anon                           0 
D u64          nr_isolated_file                           0 
D u64          workingset_nodes                           0 
D u64          workingset_refault_anon                    0 
D u64          workingset_refault_file                    0 
D u64          workingset_activate_anon                   0 
D u64          workingset_activate_file                   0 
D u64          workingset_restore_anon                    0 
D u64          workingset_restore_file                    0 
D u64          workingset_nodereclaim                     0 
D u64          nr_anon_pages                              0 
D u64          nr_mapped                                  0 
D u64          nr_file_pages                              0 
D u64          nr_dirty                                   0 
D u64          nr_writeback                               0 
D u64          nr_writeback_temp                          0 
D u64          nr_shmem                                   0 
D u64          nr_shmem_hugepages                         0 
D u64          nr_shmem_pmdmapped                         0 
D u64          nr_file_hugepages                          0 
D u64          nr_file_pmdmapped                          0 
D u64          nr_anon_transparent_hugepages              0 
D u64          nr_vmscan_write                            0 
D u64          nr_vmscan_immediate_reclaim                0 
D u64          nr_dirtied                                 0 
D u64          nr_written                                 0 
D u64          nr_kernel_misc_reclaimable                 0 
D u64          nr_foll_pin_acquired                       0 
D u64          nr_foll_pin_released                       0 
D u64          nr_kernel_stack                            0 
D u64          nr_page_table_pages                        0 
D u64          nr_swapcached                              0 
D u64          nr_dirty_threshold                         0 
D u64          nr_dirty_background_threshold              0 
D u64          pgpgin                                     0 
D u64          pgpgout                                    0 
D u64          pswpin                                     0 
D u64          pswpout                                    0 
D u64          pgalloc_dma                                0 
D u64          pgalloc_dma32                              0 
D u64          pgalloc_normal                             0 
D u64          pgalloc_movable                            0 
D u64          allocstall_dma                             0 
D u64          allocstall_dma32                           0 
D u64          allocstall_normal                          0 
D u64          allocstall_movable                         0 
D u64          pgskip_dma                                 0 
D u64          pgskip_dma32                               0 
D u64          pgskip_normal                              0 
D u64          pgskip_movable                             0 
D u64          pgfree                                     0 
D u64          pgactivate                                 0 
D u64          pgdeactivate                               0 
D u64          pglazyfree                                 0 
D u64          pgfault                                    0 
D u64          pgmajfault                                 0 
D u64          pglazyfreed                                0 
D u64          pgrefill                                   0 
D u64          pgreuse                                    0 
D u64          pgsteal_kswapd                             0 
D u64          pgsteal_direct                             0 
D u64          pgdemote_kswapd                            0 
D u64          pgdemote_direct                            0 
D u64          pgscan_kswapd                              0 
D u64          pgscan_direct                              0 
D u64          pgscan_direct_throttle                     0 
D u64          pgscan_anon                                0 
D u64          pgscan_file                                0 
D u64          pgsteal_anon                               0 
D u64          pgsteal_file                               0 
D u64          zone_reclaim_failed                        0 
D u64          pginodesteal                               0 
D u64          slabs_scanned                              0 
D u64          kswapd_inodesteal                          0 
D u64          kswapd_low_wmark_hit_quickly               0 
D u64          kswapd_high_wmark_hit_quickly              0 
D u64          pageoutrun                                 0 
D u64          pgrotated                                  0 
D u64          drop_pagecache                             0 
D u64          drop_slab                                  0 
D u64          oom_kill                                   0 
D u64          numa_pte_updates                           0 
D u64          numa_huge_pte_updates                      0 
D u64          numa_hint_faults                           0 
D u64          numa_hint_faults_local                     0 
D u64          numa_pages_migrated                        0 
D u64          pgmigrate_success                          0 
D u64          pgmigrate_fail                             0 
D u64          thp_migration_success                      0 
D u64          thp_migration_fail                         0 
D u64          thp_migration_split                        0 
D u64          compact_migrate_scanned                    0 
D u64          compact_free_scanned                       0 
D u64          compact_isolated                           0 
D u64          compact_stall                              0 
D u64          compact_fail                               0 
D u64          compact_success                            0 
D u64          compact_daemon_wake                        0 
D u64          compact_daemon_migrate_scanned             0 
D u64          compact_daemon_free_scanned                0 
D u64          htlb_buddy_alloc_success                   0 
D u64          htlb_buddy_alloc_fail                      0 
D u64          unevictable_pgs_culled                     0 
D u64          unevictable_pgs_scanned                    0 
D u64          unevictable_pgs_rescued                    0 
D u64          unevictable_pgs_mlocked                    0 
D u64          unevictable_pgs_munlocked                  0 
D u64          unevictable_pgs_cleared                    0 
D u64          unevictable_pgs_stranded                   0 
D u64          thp_fault_alloc                            0 
D u64          thp_fault_fallback                         0 
D u64          thp_fault_fallback_charge                  0 
D u64          thp_collapse_alloc                         0 
D u64          thp_collapse_alloc_failed                  0 
D u64          thp_file_alloc                             0 
D u64          thp_file_fallback                          0 
D u64          thp_file_fallback_charge                   0 
D u64          thp_file_mapped                            0 
D u64          thp_split_page                             0 
D u64          thp_split_page_failed                      0 
D u64          thp_deferred_split_page                    0 
D u64          thp_split_pmd                              0 
D u64          thp_split_pud                              0 
D u64          thp_zero_page_alloc                        0 
D u64          thp_zero_page_alloc_failed                 0 
D u64          thp_swpout                                 0 
D u64          thp_swpout_fallback                        0 
D u64          balloon_inflate                            0 
D u64          balloon_deflate                            0 
D u64          balloon_migrate                            0 
D u64          swap_ra                                    0 
D u64          swap_ra_hit                                0 
D u64          direct_map_level2_splits                   0 
D u64          direct_map_level3_splits                   0 
D u64          nr_unstable                                0 

METRICS OPERATOR TIMEPOINT
Schema         Instance                 Flags  Msize  Dsize  Hsize  UID    GID    Perm       Update            Duration          Info    
-------------- ------------------------ ------ ------ ------ ------ ------ ------ ---------- ----------------- ----------------- --------
vmstat         metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/vmstat    CL    8920   1384      0      0      0 -rwxr-xr-x 1692230596.051634          0.000128 "updt_hint_us"="1000000:50000" 
meminfo        metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/meminfo    CL    2888    520      0  12345  12345 -rwxr-xr-x 1692230596.051503          0.000069 "updt_hint_us"="1000000:50000" 
-------------- ------------------------ ------ ------ ------ ------ ------ ------ ---------- ----------------- ----------------- --------
Total Sets: 2, Meta Data (kB): 11.81, Data (kB) 1.90, Memory (kB): 13.71

=======================================================================

metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/meminfo: consistent, last update: Thu Aug 17 00:03:16 2023 +0000 [51503us] 
M u64          component_id                               1 
D u64          job_id                                     0 
D u64          app_id                                     0 
D u64          MemTotal                                   16051116 
D u64          MemFree                                    411072 
D u64          MemAvailable                               2392884 
D u64          Buffers                                    173168 
D u64          Cached                                     4036536 
D u64          SwapCached                                 91056 
D u64          Active                                     3174264 
D u64          Inactive                                   10883828 
D u64          Active(anon)                               2081720 
D u64          Inactive(anon)                             10093376 
D u64          Active(file)                               1092544 
D u64          Inactive(file)                             790452 
D u64          Unevictable                                531496 
D u64          Mlocked                                    432 
D u64          SwapTotal                                  2097148 
D u64          SwapFree                                   60 
D u64          Dirty                                      51720 
D u64          Writeback                                  4 
D u64          AnonPages                                  10289164 
D u64          Mapped                                     883272 
D u64          Shmem                                      2327760 
D u64          KReclaimable                               441432 
D u64          Slab                                       733876 
D u64          SReclaimable                               441432 
D u64          SUnreclaim                                 292444 
D u64          KernelStack                                41632 
D u64          PageTables                                 113556 
D u64          NFS_Unstable                               0 
D u64          Bounce                                     0 
D u64          WritebackTmp                               0 
D u64          CommitLimit                                10122704 
D u64          Committed_AS                               30676192 
D u64          VmallocTotal                               34359738367 
D u64          VmallocUsed                                81628 
D u64          VmallocChunk                               0 
D u64          Percpu                                     16384 
D u64          HardwareCorrupted                          0 
D u64          AnonHugePages                              18432 
D u64          ShmemHugePages                             0 
D u64          ShmemPmdMapped                             0 
D u64          FileHugePages                              0 
D u64          FilePmdMapped                              0 
D u64          HugePages_Total                            0 
D u64          HugePages_Free                             0 
D u64          HugePages_Rsvd                             0 
D u64          HugePages_Surp                             0 
D u64          Hugepagesize                               2048 
D u64          Hugetlb                                    0 
D u64          DirectMap4k                                2749276 
D u64          DirectMap2M                                10561536 
D u64          DirectMap1G                                4194304 

metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/vmstat: consistent, last update: Thu Aug 17 00:03:16 2023 +0000 [51634us] 
M u64          component_id                               1 
D u64          job_id                                     0 
D u64          app_id                                     0 
D u64          nr_free_pages                              102768 
D u64          nr_zone_inactive_anon                      2523409 
D u64          nr_zone_active_anon                        520430 
D u64          nr_zone_inactive_file                      197613 
D u64          nr_zone_active_file                        273136 
D u64          nr_zone_unevictable                        132874 
D u64          nr_zone_write_pending                      12931 
D u64          nr_mlock                                   108 
D u64          nr_bounce                                  0 
D u64          nr_zspages                                 0 
D u64          nr_free_cma                                0 
D u64          numa_hit                                   12662065045 
D u64          numa_miss                                  0 
D u64          numa_foreign                               0 
D u64          numa_interleave                            3052 
D u64          numa_local                                 12662046633 
D u64          numa_other                                 0 
D u64          nr_inactive_anon                           2523344 
D u64          nr_active_anon                             520430 
D u64          nr_inactive_file                           197613 
D u64          nr_active_file                             273136 
D u64          nr_unevictable                             132874 
D u64          nr_slab_reclaimable                        110358 
D u64          nr_slab_unreclaimable                      73111 
D u64          nr_isolated_anon                           0 
D u64          nr_isolated_file                           0 
D u64          workingset_nodes                           34917 
D u64          workingset_refault_anon                    900609 
D u64          workingset_refault_file                    232700013 
D u64          workingset_activate_anon                   65073 
D u64          workingset_activate_file                   74541615 
D u64          workingset_restore_anon                    6296 
D u64          workingset_restore_file                    45586725 
D u64          workingset_nodereclaim                     3746802 
D u64          nr_anon_pages                              2572291 
D u64          nr_mapped                                  220818 
D u64          nr_file_pages                              1075190 
D u64          nr_dirty                                   12930 
D u64          nr_writeback                               1 
D u64          nr_writeback_temp                          0 
D u64          nr_shmem                                   581940 
D u64          nr_shmem_hugepages                         0 
D u64          nr_shmem_pmdmapped                         0 
D u64          nr_file_hugepages                          0 
D u64          nr_file_pmdmapped                          0 
D u64          nr_anon_transparent_hugepages              9 
D u64          nr_vmscan_write                            2738291 
D u64          nr_vmscan_immediate_reclaim                549005 
D u64          nr_dirtied                                 269993155 
D u64          nr_written                                 260959409 
D u64          nr_kernel_misc_reclaimable                 0 
D u64          nr_foll_pin_acquired                       329043326 
D u64          nr_foll_pin_released                       329043326 
D u64          nr_kernel_stack                            41632 
D u64          nr_page_table_pages                        28389 
D u64          nr_swapcached                              22764 
D u64          nr_dirty_threshold                         105994 
D u64          nr_dirty_background_threshold              52932 
D u64          pgpgin                                     1214106483 
D u64          pgpgout                                    1170851529 
D u64          pswpin                                     900609 
D u64          pswpout                                    2750939 
D u64          pgalloc_dma                                1 
D u64          pgalloc_dma32                              619377373 
D u64          pgalloc_normal                             12042975623 
D u64          pgalloc_movable                            0 
D u64          allocstall_dma                             0 
D u64          allocstall_dma32                           0 
D u64          allocstall_normal                          380760 
D u64          allocstall_movable                         1045948 
D u64          pgskip_dma                                 0 
D u64          pgskip_dma32                               0 
D u64          pgskip_normal                              474355403 
D u64          pgskip_movable                             0 
D u64          pgfree                                     13139365226 
D u64          pgactivate                                 275410831 
D u64          pgdeactivate                               235343678 
D u64          pglazyfree                                 623494 
D u64          pgfault                                    11841715157 
D u64          pgmajfault                                 8821132 
D u64          pglazyfreed                                165160 
D u64          pgrefill                                   508213743 
D u64          pgreuse                                    989371552 
D u64          pgsteal_kswapd                             325052168 
D u64          pgsteal_direct                             78556990 
D u64          pgdemote_kswapd                            0 
D u64          pgdemote_direct                            0 
D u64          pgscan_kswapd                              697969029 
D u64          pgscan_direct                              123977398 
D u64          pgscan_direct_throttle                     0 
D u64          pgscan_anon                                19724061 
D u64          pgscan_file                                803152815 
D u64          pgsteal_anon                               2821165 
D u64          pgsteal_file                               400930138 
D u64          zone_reclaim_failed                        0 
D u64          pginodesteal                               18237 
D u64          slabs_scanned                              2874034820 
D u64          kswapd_inodesteal                          96641778 
D u64          kswapd_low_wmark_hit_quickly               54787 
D u64          kswapd_high_wmark_hit_quickly              17027 
D u64          pageoutrun                                 123624 
D u64          pgrotated                                  3210852 
D u64          drop_pagecache                             0 
D u64          drop_slab                                  0 
D u64          oom_kill                                   10 
D u64          numa_pte_updates                           0 
D u64          numa_huge_pte_updates                      0 
D u64          numa_hint_faults                           0 
D u64          numa_hint_faults_local                     0 
D u64          numa_pages_migrated                        0 
D u64          pgmigrate_success                          472793405 
D u64          pgmigrate_fail                             3526993741 
D u64          thp_migration_success                      0 
D u64          thp_migration_fail                         0 
D u64          thp_migration_split                        0 
D u64          compact_migrate_scanned                    8249426104 
D u64          compact_free_scanned                       32912179625 
D u64          compact_isolated                           4483381635 
D u64          compact_stall                              69722 
D u64          compact_fail                               69597 
D u64          compact_success                            125 
D u64          compact_daemon_wake                        74650 
D u64          compact_daemon_migrate_scanned             41663618 
D u64          compact_daemon_free_scanned                283723703 
D u64          htlb_buddy_alloc_success                   0 
D u64          htlb_buddy_alloc_fail                      0 
D u64          unevictable_pgs_culled                     1050675304 
D u64          unevictable_pgs_scanned                    1134600478 
D u64          unevictable_pgs_rescued                    1050487529 
D u64          unevictable_pgs_mlocked                    34575 
D u64          unevictable_pgs_munlocked                  21786 
D u64          unevictable_pgs_cleared                    5676 
D u64          unevictable_pgs_stranded                   5629 
D u64          thp_fault_alloc                            13822 
D u64          thp_fault_fallback                         14523 
D u64          thp_fault_fallback_charge                  0 
D u64          thp_collapse_alloc                         2031 
D u64          thp_collapse_alloc_failed                  13785 
D u64          thp_file_alloc                             0 
D u64          thp_file_fallback                          0 
D u64          thp_file_fallback_charge                   0 
D u64          thp_file_mapped                            0 
D u64          thp_split_page                             73 
D u64          thp_split_page_failed                      0 
D u64          thp_deferred_split_page                    951 
D u64          thp_split_pmd                              3626 
D u64          thp_split_pud                              0 
D u64          thp_zero_page_alloc                        25 
D u64          thp_zero_page_alloc_failed                 4240 
D u64          thp_swpout                                 0 
D u64          thp_swpout_fallback                        1 
D u64          balloon_inflate                            0 
D u64          balloon_deflate                            0 
D u64          balloon_migrate                            0 
D u64          swap_ra                                    330648 
D u64          swap_ra_hit                                181316 
D u64          direct_map_level2_splits                   1332 
D u64          direct_map_level3_splits                   9 
D u64          nr_unstable                                0 

METRICS OPERATOR TIMEPOINT
Schema         Instance                 Flags  Msize  Dsize  Hsize  UID    GID    Perm       Update            Duration          Info    
-------------- ------------------------ ------ ------ ------ ------ ------ ------ ---------- ----------------- ----------------- --------
vmstat         metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/vmstat    CL    8920   1384      0      0      0 -rwxr-xr-x 1692230606.051783          0.000255 "updt_hint_us"="1000000:50000" 
meminfo        metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/meminfo    CL    2888    520      0  12345  12345 -rwxr-xr-x 1692230606.051524          0.000112 "updt_hint_us"="1000000:50000" 
-------------- ------------------------ ------ ------ ------ ------ ------ ------ ---------- ----------------- ----------------- --------
Total Sets: 2, Meta Data (kB): 11.81, Data (kB) 1.90, Memory (kB): 13.71

=======================================================================

metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/meminfo: consistent, last update: Thu Aug 17 00:03:26 2023 +0000 [51524us] 
M u64          component_id                               1 
D u64          job_id                                     0 
D u64          app_id                                     0 
D u64          MemTotal                                   16051116 
D u64          MemFree                                    410572 
D u64          MemAvailable                               2407596 
D u64          Buffers                                    173412 
D u64          Cached                                     4021448 
D u64          SwapCached                                 91056 
D u64          Active                                     3190336 
D u64          Inactive                                   10885248 
D u64          Active(anon)                               2081736 
D u64          Inactive(anon)                             10095692 
D u64          Active(file)                               1108600 
D u64          Inactive(file)                             789556 
D u64          Unevictable                                504500 
D u64          Mlocked                                    432 
D u64          SwapTotal                                  2097148 
D u64          SwapFree                                   60 
D u64          Dirty                                      48364 
D u64          Writeback                                  4 
D u64          AnonPages                                  10294732 
D u64          Mapped                                     884284 
D u64          Shmem                                      2297748 
D u64          KReclaimable                               441484 
D u64          Slab                                       733840 
D u64          SReclaimable                               441484 
D u64          SUnreclaim                                 292356 
D u64          KernelStack                                42096 
D u64          PageTables                                 114152 
D u64          NFS_Unstable                               0 
D u64          Bounce                                     0 
D u64          WritebackTmp                               0 
D u64          CommitLimit                                10122704 
D u64          Committed_AS                               30159452 
D u64          VmallocTotal                               34359738367 
D u64          VmallocUsed                                82044 
D u64          VmallocChunk                               0 
D u64          Percpu                                     16384 
D u64          HardwareCorrupted                          0 
D u64          AnonHugePages                              18432 
D u64          ShmemHugePages                             0 
D u64          ShmemPmdMapped                             0 
D u64          FileHugePages                              0 
D u64          FilePmdMapped                              0 
D u64          HugePages_Total                            0 
D u64          HugePages_Free                             0 
D u64          HugePages_Rsvd                             0 
D u64          HugePages_Surp                             0 
D u64          Hugepagesize                               2048 
D u64          Hugetlb                                    0 
D u64          DirectMap4k                                2749276 
D u64          DirectMap2M                                10561536 
D u64          DirectMap1G                                4194304 

metricset-sample-m-0-dm6ph.ms.default.svc.cluster.local/vmstat: consistent, last update: Thu Aug 17 00:03:26 2023 +0000 [51783us] 
M u64          component_id                               1 
D u64          job_id                                     0 
D u64          app_id                                     0 
D u64          nr_free_pages                              102643 
D u64          nr_zone_inactive_anon                      2523923 
D u64          nr_zone_active_anon                        520434 
D u64          nr_zone_inactive_file                      197389 
D u64          nr_zone_active_file                        277150 
D u64          nr_zone_unevictable                        126125 
D u64          nr_zone_write_pending                      12092 
D u64          nr_mlock                                   108 
D u64          nr_bounce                                  0 
D u64          nr_zspages                                 0 
D u64          nr_free_cma                                0 
D u64          numa_hit                                   12662139090 
D u64          numa_miss                                  0 
D u64          numa_foreign                               0 
D u64          numa_interleave                            3052 
D u64          numa_local                                 12662120678 
D u64          numa_other                                 0 
D u64          nr_inactive_anon                           2523923 
D u64          nr_active_anon                             520434 
D u64          nr_inactive_file                           197389 
D u64          nr_active_file                             277150 
D u64          nr_unevictable                             126125 
D u64          nr_slab_reclaimable                        110371 
D u64          nr_slab_unreclaimable                      73089 
D u64          nr_isolated_anon                           0 
D u64          nr_isolated_file                           0 
D u64          workingset_nodes                           34863 
D u64          workingset_refault_anon                    900609 
D u64          workingset_refault_file                    232703700 
D u64          workingset_activate_anon                   65073 
D u64          workingset_activate_file                   74545239 
D u64          workingset_restore_anon                    6296 
D u64          workingset_restore_file                    45586870 
D u64          workingset_nodereclaim                     3746802 
D u64          nr_anon_pages                              2573683 
D u64          nr_mapped                                  221071 
D u64          nr_file_pages                              1071479 
D u64          nr_dirty                                   12091 
D u64          nr_writeback                               1 
D u64          nr_writeback_temp                          0 
D u64          nr_shmem                                   574437 
D u64          nr_shmem_hugepages                         0 
D u64          nr_shmem_pmdmapped                         0 
D u64          nr_file_hugepages                          0 
D u64          nr_file_pmdmapped                          0 
D u64          nr_anon_transparent_hugepages              9 
D u64          nr_vmscan_write                            2738291 
D u64          nr_vmscan_immediate_reclaim                549005 
D u64          nr_dirtied                                 269993878 
D u64          nr_written                                 260960971 
D u64          nr_kernel_misc_reclaimable                 0 
D u64          nr_foll_pin_acquired                       329043326 
D u64          nr_foll_pin_released                       329043326 
D u64          nr_kernel_stack                            42096 
D u64          nr_page_table_pages                        28538 
D u64          nr_swapcached                              22764 
D u64          nr_dirty_threshold                         106726 
D u64          nr_dirty_background_threshold              53298 
D u64          pgpgin                                     1214114941 
D u64          pgpgout                                    1170858993 
D u64          pswpin                                     900609 
D u64          pswpout                                    2750939 
D u64          pgalloc_dma                                1 
D u64          pgalloc_dma32                              619377373 
D u64          pgalloc_normal                             12043049668 
D u64          pgalloc_movable                            0 
D u64          allocstall_dma                             0 
D u64          allocstall_dma32                           0 
D u64          allocstall_normal                          380760 
D u64          allocstall_movable                         1045948 
D u64          pgskip_dma                                 0 
D u64          pgskip_dma32                               0 
D u64          pgskip_normal                              474355403 
D u64          pgskip_movable                             0 
D u64          pgfree                                     13139441479 
D u64          pgactivate                                 275411589 
D u64          pgdeactivate                               235343678 
D u64          pglazyfree                                 623494 
D u64          pgfault                                    11841785904 
D u64          pgmajfault                                 8821244 
D u64          pglazyfreed                                165160 
D u64          pgrefill                                   508213743 
D u64          pgreuse                                    989383314 
D u64          pgsteal_kswapd                             325052168 
D u64          pgsteal_direct                             78556990 
D u64          pgdemote_kswapd                            0 
D u64          pgdemote_direct                            0 
D u64          pgscan_kswapd                              697969029 
D u64          pgscan_direct                              123977398 
D u64          pgscan_direct_throttle                     0 
D u64          pgscan_anon                                19724061 
D u64          pgscan_file                                803152815 
D u64          pgsteal_anon                               2821165 
D u64          pgsteal_file                               400930138 
D u64          zone_reclaim_failed                        0 
D u64          pginodesteal                               18237 
D u64          slabs_scanned                              2874034820 
D u64          kswapd_inodesteal                          96641778 
D u64          kswapd_low_wmark_hit_quickly               54787 
D u64          kswapd_high_wmark_hit_quickly              17027 
D u64          pageoutrun                                 123624 
D u64          pgrotated                                  3210852 
D u64          drop_pagecache                             0 
D u64          drop_slab                                  0 
D u64          oom_kill                                   10 
D u64          numa_pte_updates                           0 
D u64          numa_huge_pte_updates                      0 
D u64          numa_hint_faults                           0 
D u64          numa_hint_faults_local                     0 
D u64          numa_pages_migrated                        0 
D u64          pgmigrate_success                          472793405 
D u64          pgmigrate_fail                             3526993741 
D u64          thp_migration_success                      0 
D u64          thp_migration_fail                         0 
D u64          thp_migration_split                        0 
D u64          compact_migrate_scanned                    8249426104 
D u64          compact_free_scanned                       32912179625 
D u64          compact_isolated                           4483381635 
D u64          compact_stall                              69722 
D u64          compact_fail                               69597 
D u64          compact_success                            125 
D u64          compact_daemon_wake                        74650 
D u64          compact_daemon_migrate_scanned             41663618 
D u64          compact_daemon_free_scanned                283723703 
D u64          htlb_buddy_alloc_success                   0 
D u64          htlb_buddy_alloc_fail                      0 
D u64          unevictable_pgs_culled                     1050687009 
D u64          unevictable_pgs_scanned                    1134619187 
D u64          unevictable_pgs_rescued                    1050505996 
D u64          unevictable_pgs_mlocked                    34575 
D u64          unevictable_pgs_munlocked                  21786 
D u64          unevictable_pgs_cleared                    5676 
D u64          unevictable_pgs_stranded                   5629 
D u64          thp_fault_alloc                            13822 
D u64          thp_fault_fallback                         14523 
D u64          thp_fault_fallback_charge                  0 
D u64          thp_collapse_alloc                         2031 
D u64          thp_collapse_alloc_failed                  13785 
D u64          thp_file_alloc                             0 
D u64          thp_file_fallback                          0 
D u64          thp_file_fallback_charge                   0 
D u64          thp_file_mapped                            0 
D u64          thp_split_page                             73 
D u64          thp_split_page_failed                      0 
D u64          thp_deferred_split_page                    951 
D u64          thp_split_pmd                              3626 
D u64          thp_split_pud                              0 
D u64          thp_zero_page_alloc                        25 
D u64          thp_zero_page_alloc_failed                 4240 
D u64          thp_swpout                                 0 
D u64          thp_swpout_fallback                        1 
D u64          balloon_inflate                            0 
D u64          balloon_deflate                            0 
D u64          balloon_migrate                            0 
D u64          swap_ra                                    330648 
D u64          swap_ra_hit                                181316 
D u64          direct_map_level2_splits                   1332 
D u64          direct_map_level3_splits                   9 
D u64          nr_unstable                                0 

METRICS OPERATOR COLLECTION END
```

When you are done, the JobsSet, associated jobs, and pods will be completed:

```bash
$ kubectl get jobset
```
```console
NAME               RESTARTS   COMPLETED   AGE
metricset-sample              True        2m28s
```
```bash
$ kubectl get jobs
```
```console
NAME                   COMPLETIONS   DURATION   AGE
metricset-sample-m-0   1/1           13s        2m51s
```
```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS      RESTARTS   AGE
metricset-sample-m-0-0-rq4q9   0/1     Completed   0          3m19s
```

When you are done, cleanup!

```bash
kubectl delete -f metrics.yaml
```
