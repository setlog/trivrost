module github.com/setlog/trivrost

require (
	git.sr.ht/~tslocum/preallocate v0.1.2
	github.com/MMulthaupt/chronometry v0.1.1
	github.com/andlabs/ui v0.0.0-20200610043537-70a69d6ae31e
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.13.0
	github.com/go-ole/go-ole v1.2.6
	github.com/gofrs/flock v0.8.1
	github.com/mattn/go-ieproxy v0.0.1
	github.com/prometheus/client_golang v1.12.2
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.4
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/net v0.0.0-20220531201128-c960675eff93
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a
)

go 1.16

replace git.sr.ht/~tslocum/preallocate => github.com/smallnest/preallocate v0.1.1
