calculateFit() { time curl http://46.101.82.87/admin/calculateFit\?stationID\=$1\&from\=1465653992\&to\=1471996800; }
testFit() { time Rscript fitTester.R $1 1472007600 1472383932; }
ctFit() { calculateFit $1 && testFit $1; }
