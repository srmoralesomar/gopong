package main

import (
	"math"
)

// ResolveCollisions handles AABB vs Circle continuous checks and responses
// using the Closest Point clamping algorithm to handle edge and corner hits.
func ResolveCollisions(ball *Ball, paddle *Paddle) {
	// 1. Find the closest point on the paddle's AABB to the ball's center.
	paddleCenterX := paddle.X + paddle.Width/2
	paddleCenterY := paddle.Y + paddle.Height/2
	
	halfWidth := paddle.Width / 2
	halfHeight := paddle.Height / 2
	
	// Vector from paddle center to ball center
	diffX := ball.X - paddleCenterX
	diffY := ball.Y - paddleCenterY
	
	// Clamp the difference vector to the paddle's half-extents
	clampedX := math.Max(-halfWidth, math.Min(halfWidth, diffX))
	clampedY := math.Max(-halfHeight, math.Min(halfHeight, diffY))
	
	// Closest point on the AABB to the circle center
	closestX := paddleCenterX + clampedX
	closestY := paddleCenterY + clampedY
	
	// Vector from closest point to ball center
	distanceX := ball.X - closestX
	distanceY := ball.Y - closestY
	
	distanceSq := distanceX*distanceX + distanceY*distanceY
	
	// 2. Check for collision
	if distanceSq <= ball.Radius*ball.Radius && distanceSq > 0 {
		distance := math.Sqrt(distanceSq)
		penetration := ball.Radius - distance
		
		// Normal vector (normalized distance vector)
		nx := distanceX / distance
		ny := distanceY / distance
		
		// 3. Positional Correction to prevent sticking/tunneling
		ball.X += nx * penetration
		ball.Y += ny * penetration
		
		// 4. Calculate Velocity Response
		// If collision is mostly horizontal (front face) or corner
		if math.Abs(nx) > math.Abs(ny) {
			ball.Vx = -ball.Vx
			
			// Relative bounce: hitting edges of the paddle gives a sharper angle
			hitRelative := (ball.Y - paddleCenterY) / halfHeight
			ball.Vy += hitRelative * 0.75 
			
			// Normalize velocity vector to maintain speed
			speed := math.Sqrt(ball.Vx*ball.Vx + ball.Vy*ball.Vy)
			ball.Vx = (ball.Vx / speed)
			ball.Vy = (ball.Vy / speed)
			
			// Ensure ball moves away from the paddle
			if ball.X > paddleCenterX && ball.Vx < 0 {
				ball.Vx = -ball.Vx
			} else if ball.X < paddleCenterX && ball.Vx > 0 {
				ball.Vx = -ball.Vx
			}

			ball.IncreaseSpeed()
		} else {
			// Hit the top or bottom of the paddle
			ball.Vy = -ball.Vy
			// Ensure it moves away vertically
			if ball.Y > paddleCenterY && ball.Vy < 0 {
				ball.Vy = -ball.Vy
			} else if ball.Y < paddleCenterY && ball.Vy > 0 {
				ball.Vy = -ball.Vy
			}
		}
	} else if distanceSq == 0 {
		// Deep penetration fallback (ball center exactly inside AABB)
		if ball.X > paddleCenterX {
			ball.X = paddle.X + paddle.Width + ball.Radius
			ball.Vx = math.Abs(ball.Vx)
		} else {
			ball.X = paddle.X - ball.Radius
			ball.Vx = -math.Abs(ball.Vx)
		}
	}
}
