module Main exposing (..)

import Html exposing (Html, button, div, program, text)
import Html.Events exposing (onClick)

type alias Model = 
    { test : String }

type Msg = Test

main : Program Never Model Msg
main =
    program { init = init, update = update, subscriptions = subscriptions, view = view } 

init : (Model, Cmd Msg)
init =
    Model "" ! []

subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
    Test ->
        { model | test = "Changed" } ! []

view : Model -> Html Msg
view model =
    div [] [ if String.isEmpty model.test then text "Hello world" else text model.test, button [ onClick Test ] [ text "Change" ] ]