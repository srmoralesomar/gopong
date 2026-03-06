# Go Pong

A classic Pong clone built with Go and the [Ebitengine](https://ebitengine.org/) 2D game library. 

This project aims to provide a robust, readable foundation for a 2D paddle-and-ball game, demonstrating advanced continuous collision detection (CCD) and a clean separation between game state, physics, and rendering logic.

## How to Play

### Prerequisites
- Go 1.24+ installed on your system.
- An environment capable of running CGO (Ebitengine relies on this for windowing/graphics). On macOS and most Linux distros, this works out-of-the-box. On Windows, a C compiler like GCC (via MSYS2 or MinGW) might be necessary depending on your setup.

### Running the game
Clone the repository and run:
```bash
go run .
```

### Controls
- **Up Arrow / Down Arrow**: Move the left paddle.
- **Spacebar**: Serve the ball at the start of the game or after a point is scored.

## Architecture & Ebitengine

The game is powered by the **Ebitengine** library. Ebitengine is a "dead simple" 2D game engine for Go. It operates on a very simple, predictable game loop interface that the `Game` struct (in `game.go`) implements:

1.  **`Update()`**: Called 60 times per second (by default). This is where the game state is modified. We handle keyboard input (`ebiten.IsKeyPressed`), move the paddles, update the ball's position, and check for collisions here. No drawing happens in `Update()`.
2.  **`Draw(screen *ebiten.Image)`**: Called every frame to render the game out to the screen. We clear the screen, and then instruct the paddles, ball, and scoreboard to draw themselves onto the `screen` image object using Ebiten's `vector` and `text/v2` sub-packages.
3.  **`Layout(outsideWidth, outsideHeight int)`**: Determines the logical screen size and how it scales to the physical window. We use a fixed `800x600` resolution.

Because Ebitengine handles the complex platform-specific rendering (OpenGL, Metal, DirectX, etc.) and window management, the code remains highly focused on the actual game logic.

## Collision Detection Engine

The core of this implementation is the **robust physics engine** found in `physics.go`. 

In many introductory game tutorials, paddle collisions are handled using simple AABB (Axis-Aligned Bounding Box) intersection checks, treating the ball as a square. This often leads to "tunneling" (where a fast ball skips past a thin paddle entirely between frames) or "sticking" (where the ball gets trapped inside the paddle and rapidly oscillates back and forth).

This project solves this by using a **Circle vs AABB Closest Point Algorithm** combined with **Positional Correction**.

### 1. Closest Point Algorithm
When checking for a collision between the circular ball and the rectangular paddle:
1. We calculate the vector distance between the ball's center and the paddle's center.
2. We **clamp** that distance vector so its X/Y coordinates do not exceed the paddle's actual half-width and half-height limits.
3. This gives us the exact coordinates of the closest point on the paddle's surface to the ball's center.
4. If the distance from the ball's center to this localized "closest point" is less than or equal to the ball's radius, a collision has occurred!

By using this math, we accurately detect if the ball hit the flat front face, the top/bottom edges, or exactly on an angled corner of the paddle.

### 2. Positional Correction
When a hit is registered, the ball has usually already penetrated slightly *inside* the paddle's boundaries during that 1/60th of a second frame update. 
Before reversing the velocity, the engine calculates the exact penetration depth and immediately forces the ball's X/Y coordinates backwards along the collision normal so it sits perfectly flush with the edge of the paddle. 

**This makes it physically impossible for the ball to ever get stuck inside walls or paddles.**

### 3. Bounce Physics (Relative Velocity)
In Pong, if the ball hits the absolute edge of a paddle, it should bounce away at a sharper vertical angle than if it hits dead-center. 
When a front-face collision happens, the algorithm checks how far away the ball's Y-coordinate is from the paddle's center Y-coordinate. It uses this relative distance to dynamically increase or decrease the ball's `Vy` (Vertical Velocity), creating natural-feeling gameplay that allows players to aim their shots.
