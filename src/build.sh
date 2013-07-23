export GOPATH=`pwd`/../:`pwd`:$GOPATH
go env
go get github.com/cihub/seelog
go get github.com/zeebo/sbloom
go get github.com/robfig/config
go get github.com/PuerkitoBio/purell
make build

