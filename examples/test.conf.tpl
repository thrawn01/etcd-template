# first create a key with our values in it
# $ etcdctl set /mailgun/configs/ord/service1 '{"value1": 1, "value2": 2}'
# Now run $ etcd-template /mailgun/configs/ord/service1 examples
# Inspect the contents of the output
# $ cat examples/test.conf
key1={{ .value1 }}
key2={{ .value2 }}
key3=3
