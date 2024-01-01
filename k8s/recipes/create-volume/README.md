kubectl cp ./hello.wasm simplism/task-pv-pod:data/hello.wasm

#kubectl exec -n simplism -it task-pv-pod -- /bin/bash
cd data
ls

kubectl delete -f volume.yaml -n simplism

kubectl exec -n simplism -it task-pv-pod -- /bin/sh
cd data
ls
