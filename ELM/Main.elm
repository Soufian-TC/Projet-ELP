module Main exposing (..)
import Drawing exposing (display, Program)
import Browser
import Html exposing (Html, div, input, text, button)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)

-- MAIN
main =
    Browser.sandbox { init = init, update = update, view = view }


-- MODEL
type alias Model =
    { content : String
    , instructions : List (String, Int) -- Liste des instructions avec leur répétition
    , showDrawing : Bool
    ,errorMessage : Maybe String
    }


init : Model
init =
    { content = ""
    , instructions = []
    , showDrawing = False
    , errorMessage = Nothing
    }


-- UPDATE
type Msg
    = UpdateInput String -- Met à jour le contenu tapé
    | ParseInstructions  -- Parse les instructions après validation
    | ToggleDrawing       -- Affiche ou masque la fenêtre de dessin


update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newContent ->
            { model | content = newContent }

        ParseInstructions ->
            let
                parsedInstructions = parseInstructions model.content
            in
                case parsedInstructions of
                Just instr -> { model | instructions = instr, errorMessage = Nothing }
                Nothing -> { model | errorMessage = Just "Erreur : Vérifiez votre syntaxe !" }
        ToggleDrawing ->
            { model | showDrawing = not model.showDrawing }


-- PARSER
parseInstructions : String -> Maybe (List (String, Int))
parseInstructions input =
    let
        -- Fonction pour couper la chaîne d'entrée en tokens
        tokenize : String -> List String
        tokenize str =
            String.split " " str -- On sépare les mots par des espaces

        -- Fonction pour parser une instruction imbriquée
        parseInstruction : List String -> (List (String, Int), List String)  -- 🔴 Déplacement avant son utilisation
        parseInstruction tokens =
            case tokens of
                -- Cas de crochet ouvrant, début d'un bloc
                "[" :: rest ->
                    let
                        (instructionsInside, restAfterBlock) = parseBlock rest
                    in
                    (instructionsInside, restAfterBlock)

                -- Cas de crochet fermant, fin du bloc
                "]" :: rest ->
                    ([], rest)

                -- Cas de l'instruction Forward
                "Forward" :: value :: rest ->
                    case String.toInt value of
                        Just n ->
                            let
                                (nextInstructions, remainingTokens) = parseInstruction rest
                            in
                            ([("Forward", n)] ++ nextInstructions, remainingTokens)
                        Nothing ->
                            parseInstruction rest

                -- Cas de l'instruction Right
                "Right" :: value :: rest ->
                    case String.toInt value of
                        Just n ->
                            let
                                (nextInstructions, remainingTokens) = parseInstruction rest
                            in
                            ([("Right", n)] ++ nextInstructions, remainingTokens)
                        Nothing ->
                            parseInstruction rest

                -- Cas de l'instruction Left
                "Left" :: value :: rest ->
                    case String.toInt value of
                        Just n ->
                            let
                                (nextInstructions, remainingTokens) = parseInstruction rest
                            in
                            ([("Left", n)] ++ nextInstructions, remainingTokens)
                        Nothing ->
                            parseInstruction rest

                -- Cas du Repeat
                "Repeat" :: nStr :: "[" :: rest ->
                    case String.toInt nStr of
                        Just n ->
                            let
                                (instructionsInside, restAfterRepeat) = parseBlock rest
                                repeatedInstructions = List.concat (List.repeat n instructionsInside)
                                (nextInstructions, remainingTokens) = parseInstruction restAfterRepeat
                            in
                            (repeatedInstructions ++ nextInstructions, remainingTokens)
                        Nothing ->
                            parseInstruction rest

                -- Si on rencontre une virgule, on l'ignore et on continue avec les tokens suivants
                "," :: rest ->
                    parseInstruction rest

                -- Sinon, on termine ou on passe à la prochaine instruction
                _ ->
                    ([], [])

        -- Fonction pour parser un bloc de code entre crochets
        parseBlock : List String -> (List (String, Int), List String)
        parseBlock tokens =
            case tokens of
                -- Si on rencontre un crochet fermant, c'est la fin du bloc
                "]" :: rest ->
                    ([], rest)

                -- Sinon, on garde les instructions et on continue
                _ ->
                    let
                        (innerInstructions, remainingTokens) = parseInstruction tokens
                    in
                    (innerInstructions, remainingTokens)

    in
    -- Retourne `Nothing` si la fonction `parseInstruction` échoue à parser tout le contenu
    case tokenize input |> parseInstruction of
        ([], []) -> Nothing  -- Si aucun résultat, c'est une erreur
        (instructions, []) -> Just instructions  -- Si succès, retourne les instructions
        _ -> Nothing  -- Si la structure des tokens ne correspond pas à ce qu'on attend
-- VIEW
view : Model -> Html Msg
view model =
    div []
        [ div [] -- Premier div
            [ text "⚠️ Veuillez respecter les espaces entre les éléments !" ]
        , div [] -- Deuxième div avec champ de saisie plus large
            [ input
                [ placeholder "Example: [ Forward 100 , Repeat 4 [ Forward 50 , Left 90 ] , Forward 100 ]"
                , value model.content
                , onInput UpdateInput
                , style "width" "100%"
                ]
                []
            ]
        , div []
            [ button [ onClick ParseInstructions ] [ text "Parse Instructions" ]
            , button [ onClick ToggleDrawing ] [ text "Toggle Drawing" ]
            ]
        , case model.errorMessage of
            Just errMsg -> div [ style "color" "red" ] [ text errMsg ]  -- 🔴 Affichage de l'erreur en rouge
            Nothing -> text ""

        , if model.showDrawing then
            div
                [ style "border" "1px solid black"
                , style "width" "500px"
                , style "height" "500px"
                ]
                [ display model.instructions ]
          else
            text ""

        , div []
            [ text "Parsed Instructions: "
            , text (String.join ", " (List.map (\(cmd, count) -> cmd ++ " " ++ String.fromInt count) model.instructions))
            ]
        ]
