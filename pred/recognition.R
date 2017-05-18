library(fftw)
library(randomForest)
library(RPostgreSQL)
library(seewave)
library(tuneR)

humanFrequency <- 280

analyzeWav <- function(path, start = 0, end = Inf) {
  file <- tempfile('tmp.wav')
  wave <- download.file(path, file, mode = 'wb')
  tuneWave <- readWave(file, from = start, to = end, units = "seconds")
  waveSpec <- spec(tuneWave, f = tuneWave@samp.rate, plot = F)
  analysis <- specprop(waveSpec, f = tuneWave@samp.rate, flim = c(0, humanFrequency / 1000), plot = F)
  
  meanfreq <- analysis$mean / 1000
  sd <- analysis$sd / 1000
  median <- analysis$median / 1000
  Q25 <- analysis$Q25 / 1000
  Q75 <- analysis$Q75 / 1000
  IQR <- analysis$IQR / 1000
  skew <- analysis$skewness
  kurt <- analysis$kurtosis
  sp.ent <- analysis$sh
  sfm <- analysis$sfm
  mode <- analysis$mode / 1000
  centroid <- analysis$cent / 1000
  
  fundamental <- fund(tuneWave, f = tuneWave@samp.rate, ovlp = 50, threshold = 5, wl = 2048,
                      ylim = c(0, humanFrequency / 1000), fmax = humanFrequency, plot = F)
  
  meanfun <- mean(fundamental[, 'y'], na.rm = T)
  minfun <- min(fundamental[, 'y'], na.rm = T)
  maxfun <- max(fundamental[, 'y'], na.rm = T)
  
  b <- c(0, 22)
  dom <- dfreq(tuneWave, f = tuneWave@samp.rate, wl = 2048, ylim = c(0, humanFrequency / 1000),
               ovlp = 0, threshold = 5, bandpass = b * 1000, fftw = T, plot = F)[, 2]
  
  meandom <- mean(dom, na.rm = TRUE)
  mindom <- min(dom, na.rm = TRUE)
  maxdom <- max(dom, na.rm = TRUE)
  dfrange <- (maxdom - mindom)
  duration <- (end - start)
  
  changes <- vector()
  for(d in which(!is.na(dom))) {
    change <- abs(dom[d] - dom[d + 1])
    changes <- append(changes, change)
  }
  if(mindom == maxdom) modindx <- 0 else modindx <- mean(changes, na.rm = T) / dfrange
  
  obj <- data.frame(duration, meanfreq, sd, median, Q25, Q75, IQR, skew, kurt, sp.ent, sfm, mode, 
                    centroid, meanfun, minfun, maxfun, meandom, mindom, maxdom, dfrange, modindx)
  names(obj) <- c("duration", "meanfreq", "sd", "median", "Q25", "Q75", "IQR", "skew", "kurt", "sp.ent", 
                  "sfm","mode", "centroid", "meanfun", "minfun", "maxfun", "meandom", "mindom", "maxdom",
                  "dfrange", "modindx")
  obj
}

sql_voice <- function(voice, pred) {
  library(RPostgreSQL)
  pg <- dbDriver("PostgreSQL")
  con <- dbConnect(drv = pg,
                   host = "ec2-54-221-255-153.compute-1.amazonaws.com",
                   port = 5432,
                   dbname = "d5lt84kg6d6kpd",
                   user = 'cpsnnurbjhrrmf',
                   password = 'bf11b782407c2085826ae94cb2d10d5f7b178660dabf77e7e14e85b6bfab77da')
  q1 <- "BEGIN; INSERT INTO voices ("
  q11 <- apply(as.matrix(names(voice)), 1, function(x) paste("\"", x, "\"", collapse = ", ", sep = ""))
  q12 <- ") VALUES"
  q2 <- apply(as.matrix(voice), 1, function(x) paste("(", paste(x, collapse=","),
                                      ")", sep=""))
  q3 <- "; COMMIT;"
  qry <- paste(q1, paste(q11, collapse = ","), q12, paste(q2, collapse=","), q3)
  dbGetQuery(con, qry)
  
  voiceId <- dbGetQuery(con, "SELECT currval('voices_id_seq');")
  q4 <- "BEGIN; INSERT INTO predictions (voice_id_fk, sex) VALUES ("
  qry2 <- paste(q4, paste(voiceId), ", ", paste("'", pred, "'", sep = ""), ") ; COMMIT;", sep = "")
  dbGetQuery(con, qry2)

  voiceId
}

path <- commandArgs(trailingOnly = T)
analyzedVoice <- analyzeWav(path, end = 20)
model.forest <- readRDS('/app/pred/model.forest.rds')
prediction <- predict(model.forest, analyzedVoice)
cat(prediction, sql_voice(analyzedVoice, prediction))