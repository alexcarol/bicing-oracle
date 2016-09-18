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
    "SELECT bikes, weather_type, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation_2 WHERE id=%d AND updatetime >= FROM_UNIXTIME(%d)",
    stationID,
    from
)
data <- dbGetQuery(mydb, query)


isCalm <- function(weather) {
    return(weather >= 3 && weather <= 8)
}

data$pbikes <- data$bikes > 1
data$weather <- isCalm(data$weather)
data$dayMoment <- data$updatetime %% 86400
data$weekday <- as.POSIXlt(as.POSIXct(data$updatetime, origin="1970-01-01"))$wday

bikeFit <- randomForest(bikes ~ updatetime + dayMoment + weekday + weather_type, data=data, importance=TRUE, ntree=100)
pBikeFit <- randomForest(pbikes ~ updatetime + dayMoment + weekday + weather_type, data=data, importance=TRUE, ntree=100)


dir.create("/tmp/station/bike", recursive=TRUE)
dir.create("/tmp/station/pbike", recursive=TRUE)

saveRDS(bikeFit, sprintf("/tmp/station/bike/%d.fit", stationID))
saveRDS(pBikeFit, sprintf("/tmp/station/pbike/%d.fit", stationID))

dbDisconnect(mydb)
