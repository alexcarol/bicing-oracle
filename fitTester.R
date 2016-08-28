library(RMySQL) # will load DBI as well
library(randomForest)
library(DBI)

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
    "SELECT pbikes, bikes, slots, weather_type, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation_2 WHERE id=%d AND updatetime >= FROM_UNIXTIME(%d) AND updatetime <= FROM_UNIXTIME(%d)",
    stationID,
    from,
    to
)
data <- dbGetQuery(mydb, query)

data$dayMoment <- data$updatetime %% 86400
data$weekday <- as.POSIXlt(as.POSIXct(data$updatetime, origin="1970-01-01"))$wday


objectID <- sprintf("/tmp/station/bike/%d.fit", stationID)
fit <- readRDS(objectID)

data$prediction <- predict(fit, data)

data$predictionOk <- data$pbikes == 1 && data$prediction > 0.5 || data$pbikes == 0 && data$prediction <= 0.5

table(data$predictionOk, data$pbikes)

write.csv(data, sprintf("station%d.csv", stationID))

dbDisconnect(mydb)