# jstat_exporter

it is a fork from [wyukawa/jstat_exporter](https://github.com/wyukawa/jstat_exporter/tree/master)

## Quick start
```
curl -LO https://raw.githubusercontent.com/xibolun/jstat_exporter/main/install.sh && bash install.sh
```

you can input param like:
```
Enter Jstat path (default is /usr/bin/jstat): $(which jstat)
Enter JMS exporter port(default is 9010): 7080
Enter target java pid: $(ps -ef | grep jira | grep -vE "pts|color|jmx_exporter" | awk '{print $2}')
```
install.sh will install jstat_exporter at /opt/jstat_exporter

## Build By Yourself
1. `make build` you will get a binary file jstat_exporter.
3. start server by `jstat_exporter`
4. you can access the metrics at `http://localhost:9010/metrics`