module Drawing exposing (display, Program, Instruction(..))

import Svg exposing (Svg, line, svg)
import Svg.Attributes exposing (width, height, viewBox, stroke, x1, y1, x2, y2)

-- Types
type alias Program = List (String, Int) -- Instructions comme tuples (action, valeur)

type alias Turtle = { x : Float, y : Float, angle : Float }

type Instruction = Forward Float | Left Float | Right Float

-- Convertit des degrés en radians
degreesToRadians : Float -> Float
degreesToRadians degrees = degrees * pi / 180

-- Affichage des instructions
display : Program -> Svg msg
display instructions = 
    let centre = { x = 250, y = 250, angle = 0 }
        step : (String, Float) -> Turtle -> (Turtle, Svg msg)
        step (action, value) turtle = case action of
                "Forward" -> let radianAngle = degreesToRadians turtle.angle
                                 newX = turtle.x + (value * cos radianAngle)
                                 newY = turtle.y - (value * sin radianAngle)
                                 lineSvg = 
                                    line [ x1 (String.fromFloat turtle.x)
                                    , y1 (String.fromFloat turtle.y)
                                    , x2 (String.fromFloat newX)
                                    , y2 (String.fromFloat newY)
                                    , stroke "black"]
                                    []
                            in ( { turtle | x = newX, y = newY }, lineSvg )

                "Left" -> ( { turtle | angle = turtle.angle - value }, svg [] [] )

                "Right" -> ( { turtle | angle = turtle.angle + value }, svg [] [] )

                _ -> ( turtle, svg [] [] )

        -- Parcourt les instructions et accumule les lignes SVG
        draw : Program -> Turtle -> List (Svg msg)
        draw instruction turtle = case instruction of
                [] -> []
                instr::q -> 
                        let 
                                (action, value) = instr
                                (newTurtle, lineSvg) = step (action, toFloat value) turtle 
                        in 
                        lineSvg :: draw q newTurtle

    in
    svg [ width "500", height "500", viewBox "0 0 500 500" ]
        (draw instructions centre)
