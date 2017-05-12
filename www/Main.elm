module Main exposing (..)

import Html exposing (Html, button, div, program, text)
import Html.Events exposing (onClick)


type alias Model =
    { test : String }


type Msg
    = Test


main : Program Never Model Msg
main =
    program { init = init, update = update, subscriptions = subscriptions, view = view }


init : ( Model, Cmd Msg )
init =
    Model "" ! []


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Test ->
            { model | test = "Changed" } ! []


view : Model -> Html Msg
view model =
    div []
        [ h1 [] [ text "Detect your sex - Using your voice" ]
        , button [ id "startButton", onClick Test ] [ text "Start recording" ]
        ]
