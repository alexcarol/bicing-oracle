library(RMySQL) # will load DBI as well
library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])
from <- as.integer(args[2])
to <- as.integer(args[3])

user <- Sys.getenv("MYSQL_RAW_DATA_USER")
dbname <- Sys.getenv("MYSQL_RAW_DATA_NAME")
password <- Sys.getenv("MYSQL_RAW_DATA_PASSWORD")
host <- Sys.getenv("MYSQL_RAW_DATA_HOST")
mydb <- dbConnect(MySQL(), user=user, dbname=dbname, password=password, host=host)

query <-  sprintf(
    "SELECT pbikes, weather_type, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation WHERE id=%d AND updatetime >= FROM_UNIXTIME(%d) AND updatetime <= FROM_UNIXTIME(%d)",
    stationID,
    from,
    to
)
data <- dbGetQuery(mydb, query)

data$dayMoment <- data$updatetime %% 86400
data$weekday <- as.POSIXlt(as.POSIXct(data$updatetime, origin="1970-01-01"))$wday
data$weekdaySunday <- data$weekday == 0
data$weekdayMonday <- data$weekday == 1
data$weekdayTuesday <- data$weekday == 2
data$weekdayWednesday <- data$weekday == 3
data$weekdayThursday <- data$weekday == 4
data$weekdayFriday <- data$weekday == 5

bikeFit <- randomForest(pbikes ~ updatetime + dayMoment + weather_type + data$weekdaySunday + data$weekdayMonday + data$weekdayTuesday + data$weekdayWednesday + data$weekdayThursday + data$weekdayFriday, data=data, importance=TRUE, ntree=100)

dir.create("/tmp/station/bike", recursive=TRUE)

saveRDS(bikeFit, sprintf("/tmp/station/bike/%d.fit", stationID))

#slots will be enabled after they have been properly tested
#dir.create("/tmp/station/slot", recursive=TRUE)
#slotFit <- randomForest(pbikes ~ updatetime + dayMoment + weekday + weather_type, data=data, importance=TRUE, ntree=100)
#saveRDS(slotFit, sprintf("/tmp/station/slot/%d.fit", stationID))

dbDisconnect(mydb)
