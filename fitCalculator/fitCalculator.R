library(RMySQL) # will load DBI as well
library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])

mydb <- dbConnect(MySQL(), user='root', dbname='bicing_oracle_raw')

##TODO use a normalized table/view with the weather as well as the pbikes
query <-  sprintf(
    "SELECT pbikes, pslots, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation WHERE id=%d", ## updatetime > X?
    stationID
)
data <- dbGetQuery(mydb, query)

fit <- randomForest(as.factor(pbikes) ~ updatetime, data=data, importance=TRUE, ntree=100)


serializedObject <- serialize(fit, NULL, ascii=T)
dbDisconnect(mydb)

objectID <- sprintf("/tmp/station/station_%d.fit", stationID)
print(objectID)
saveRDS(fit, objectID)

##TODO use a normalized table/view with the weather as well
