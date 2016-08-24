library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])
predictBikes <- as.logical(args[2])
updatetime <- as.integer(args[3])
weather <- as.integer(args[4])


objectID <- sprintf("/tmp/station/bike/%d.fit", stationID)
fit <- readRDS(objectID)

weather_type <- c(weather)
weekday <- c(as.POSIXlt(as.POSIXct(updatetime, origin="1970-01-01"))$wday)
dayMoment <- c(updatetime %% 86400)
updatetime <- c(updatetime)
object <- data.frame(updatetime, dayMoment, weekday, weather_type)


predict(fit, object)