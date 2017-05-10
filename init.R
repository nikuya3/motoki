my_packages = c("seewave", "tuneR", "fftw", "randomForest", "RPostgreSQL")

install_if_missing = function(p) {
  if (p %in% rownames(installed.packages()) == FALSE) {
    install.packages(p)
  }
}

invisible(sapply(my_packages, install_if_missing))