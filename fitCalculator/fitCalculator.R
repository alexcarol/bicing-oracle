library(RMySQL) # will load DBI as well
library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])

mydb <- dbConnect(MySQL(), user='root', dbname='bicing_oracle_raw')

query <-  sprintf(
    "SELECT pbikes, weather_type, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation WHERE id=%d", ## updatetime > X?
    stationID
)
data <- dbGetQuery(mydb, query)

data$dayMoment <- data$updatetime %% 86400
data$weekday <- as.POSIXlt(as.POSIXct(data$updatetime, origin="1970-01-01"))$wday

bikeFit <- randomForest(pbikes ~ updatetime + dayMoment + weekday + weather_type, data=data, importance=TRUE, ntree=100)


dir.create("/tmp/station/bike", recursive=TRUE)

saveRDS(bikeFit, sprintf("/tmp/station/bike/%d.fit", stationID))

#slots will be enabled after they have been properly tested
#dir.create("/tmp/station/slot", recursive=TRUE)
#slotFit <- randomForest(pbikes ~ updatetime + dayMoment + weekday + weather_type, data=data, importance=TRUE, ntree=100)
#saveRDS(slotFit, sprintf("/tmp/station/slot/%d.fit", stationID))

dbDisconnect(mydb)
