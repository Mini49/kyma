---
title: Etcd Database
type: Details
---

The Service Catalog requires an `etcd` database cluster for a production use.
It has a separate `etcd` cluster defined in the Service Catalog [etcd-stateful][sc-etcd-sub-chart] sub-chart.
The [etcd-backup-operator][etcd-backup-operator] executes the backup procedure.

## Details

This section describes the backup and restore processes of the `etcd` cluster for the Service Catalog.

### Backup

To execute the backup process, you must set the following values in the [core][core-chart-values] chart:

| Property name              | Description |
|---------------------------------------------------|---|
| **global.etcdBackup.enabled**                       | If set to `true`, the [etcd-operator][etcd-operator-chart] chart and the Service Catalog [sub-chart][sc-backup-sub-chart] installs the CronJob which executes periodically the [Etcd Backup][etcd-backup-app] application. The etcd-operator also creates the [Secret][abs-creds] with the **storage-account** and **storage-key** keys. For more information on how to configure the backup CronJob, see the [Etcd Backup][etcd-backup-app-readme] documentation. |
| **global.etcdBackup.containerName**                 | The ABS container to store the backup. |
| **etcd-operator.backupOperator.abs.storageAccount** | The name of the storage account for the Azure Blob Storage (ABS). It stores the value for the **storage-account** Secret key. |
| **etcd-operator.backupOperator.abs.storageKey**     | The key value of the storage account for the ABS. It stores the value for the **storage-key** Secret key. |

> **NOTE:** If you set the **storageAccount**, **storageKey**, and **containerName** properties, the **global.etcdBackup.enabled** must be set to `true`.

### Restore

Follow this instruction to restore an `etcd` cluster from the existing backup.

1. Export the **ABS_PATH** environment variable with the path to the last successful backup file.
```bash
export ABS_PATH=$(kubectl get cm -n kyma-system sc-recorded-etcd-backup-data -o=jsonpath='{.data.abs-backup-file-path-from-last-success}')
export BACKUP_FILE_NAME=etcd.backup
```

2. Download the backup to the local workstation. You can do it from the portal or by using [azure cli][az-cli]. Set the downloaded file path:

```bash
export BACKUP_FILE_NAME=/path/to/downloaded/file
```

3. Copy the backup file to every running Pod of the StatefulSet.

```bash
for i in {0..2};
do
kubectl cp ./$BACKUP_FILE_NAME kyma-system/core-catalog-etcd-stateful-$i:/$BACKUP_FILE_NAME
done
```

4. Restore the backup on every Pod of the StatefulSet.

```bash
for i in {0..2};
do
  remoteCommand="etcdctl snapshot restore /$BACKUP_FILE_NAME "
  remoteCommand+="--name core-catalog-etcd-stateful-$i --initial-cluster "
  remoteCommand+="core-catalog-etcd-stateful-0=https://core-catalog-etcd-stateful-0.core-catalog-etcd-stateful.kyma-system.svc.cluster.local:2380,"
  remoteCommand+="core-catalog-etcd-stateful-1=https://core-catalog-etcd-stateful-1.core-catalog-etcd-stateful.kyma-system.svc.cluster.local:2380,"
  remoteCommand+="core-catalog-etcd-stateful-2=https://core-catalog-etcd-stateful-2.core-catalog-etcd-stateful.kyma-system.svc.cluster.local:2380 "
  remoteCommand+="--initial-cluster-token etcd-cluster-1 "
  remoteCommand+="--initial-advertise-peer-urls https://core-catalog-etcd-stateful-$i.core-catalog-etcd-stateful.kyma-system.svc.cluster.local:2380"

  kubectl exec core-catalog-etcd-stateful-$i -n kyma-system -- sh -c "rm -rf core-catalog-etcd-stateful-$i.etcd"
  kubectl exec core-catalog-etcd-stateful-$i -n kyma-system -- sh -c "rm -rf /var/run/etcd/backup.etcd"
  kubectl exec core-catalog-etcd-stateful-$i -n kyma-system -- sh -c "$remoteCommand"
  kubectl exec core-catalog-etcd-stateful-$i -n kyma-system -- sh -c "mv -f core-catalog-etcd-stateful-$i.etcd /var/run/etcd/backup.etcd"
  kubectl exec core-catalog-etcd-stateful-$i -n kyma-system -- sh -c "rm $BACKUP_FILE_NAME"
done
```

5. Delete old Pods.

```bash
kubectl delete pod core-catalog-etcd-stateful-0 core-catalog-etcd-stateful-1 core-catalog-etcd-stateful-2 -n kyma-system
```

[etcd-backup-operator]:https://github.com/coreos/etcd-operator/blob/master/doc/user/walkthrough/backup-operator.md

<!-- These absolute paths should be replaced with the relative links after adding this functionality to Kyma -->
[az-cli]:https://docs.microsoft.com/en-us/cli/azure/?view=azure-cli-latest

[sc-etcd-sub-chart]:https://github.com/kyma-project/kyma/blob/master/resources/core/charts/service-catalog/charts/etcd-stateful/templates
[sc-backup-sub-chart]:https://github.com/kyma-project/kyma/blob/master/resources/core/charts/service-catalog/charts/etcd-stateful/templates/05-backup-job.yaml
[etcd-operator-chart]:https://github.com/kyma-project/kyma/blob/master/resources/core/charts/etcd-operator
[etcd-backup-operator-chart]:https://github.com/kyma-project/kyma/blob/master/resources/core/charts/etcd-operator/templates/backup-deployment.yaml
[core-chart-values]:https://github.com/kyma-project/kyma/blob/master/resources/core/values.yaml

[etcd-backup-app-readme]:https://github.com/kyma-project/kyma/blob/master/tools/etcd-backup/README.md
[etcd-backup-app]:https://github.com/kyma-project/kyma/blob/master/tools/etcd-backup

[abs-creds]:https://github.com/kyma-project/kyma/blob/master/resources/core/charts/etcd-operator/templates/etcd-backup-abs-storage-secret.yaml
