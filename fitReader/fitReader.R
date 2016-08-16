library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])
updatetime <- as.integer(args[2])
# weather <- as.integer(args[3])

print("stationID")
print(stationID)

objectID <- sprintf("/tmp/station/station_%d.fit", stationID)
print(objectID)
fit <- readRDS(objectID)
print(fit)
