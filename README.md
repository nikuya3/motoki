# Motoki
A web app that utilizes an R prediction model to recognize whether a recorded voice is female or male.
The backend is a mix of Go and R. The frontend started utilizes [Recorder.js](https://github.com/mattdiamond/Recorderjs) for voice recordings. The recordings are converted to vocal statistics on the fly, so no actual voice samples are saved on the server. There's also the possibility to rate the prediction, which could be used to improve the prediction model on a later stage.
