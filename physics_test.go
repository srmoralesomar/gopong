package main

import (
	"testing"
	"math"
)

// Helper function to test float equality with epsilon
// removed: floatEquals

func TestResolveCollisions_FrontHit(t *testing.T) {
	// Setup: Player paddle is on the left
	paddle := NewPaddle(50, 200, 15, 100)
	// Ball hits center-front of paddle from the right
	ball := NewBall(75, 250, 10)
	ball.Vx = -5 // moving left
	ball.Vy = 0

	ResolveCollisions(ball, paddle)

	// Since it hit the right side, ball Vx should now be positive
	if ball.Vx <= 0 {
		t.Errorf("expected horizontal velocity to reverse, got %f", ball.Vx)
	}

	// Positional correction: Ball X should be pushed right of the paddle + radius
	expectedX := paddle.X + paddle.Width + ball.Radius
	if ball.X < expectedX {
		t.Errorf("position not corrected. expected X >= %f, got %f", expectedX, ball.X)
	}
}

func TestResolveCollisions_TopHit(t *testing.T) {
	// Setup: paddle on the left
	paddle := NewPaddle(50, 200, 15, 100)
	// Ball hits top of paddle (X is above the paddle, Y is near the top edge)
	ball := NewBall(57, 192, 10) // Hit X=57 (center of paddle width), Y=192 (above it, overlaps 200 by 2)
	ball.Vx = 0
	ball.Vy = 5 // Moving down
	
	ResolveCollisions(ball, paddle)

	// Ball should bounce up
	if ball.Vy >= 0 {
		t.Errorf("expected vertical velocity to reverse, got %f", ball.Vy)
	}

	// Positional correction: Ball Y should be pushed above the paddle
	expectedY := paddle.Y - ball.Radius
	if ball.Y > expectedY {
		t.Errorf("position not corrected. expected Y <= %f, got %f", expectedY, ball.Y)
	}
}

func TestResolveCollisions_DeepPenetration(t *testing.T) {
	// Setup: paddle on the left
	paddle := NewPaddle(50, 200, 15, 100)
	// Ball center is exactly inside the paddle (deep tunneling/glitch) 
	ball := NewBall(57, 250, 10) // Right in the middle
	ball.Vx = -5 // Was moving left

	ResolveCollisions(ball, paddle)

	// Ball should be pushed outside the paddle to the right (since center X=57 is > center X of paddle which is 57.5? Wait, paddle center X is 50 + 7.5 = 57.5)
	// Actually 57 is < 57.5. So it will push LEFT!
	// Let's test precisely based on the paddle center logic
	expectedVx := -math.Abs(ball.Vx) // It pushed left, so velocity should be negative
	if ball.Vx != expectedVx {
		t.Errorf("Deep penetration velocity incorrect. Expected %f, got %f", expectedVx, ball.Vx)
	}

	expectedX := paddle.X - ball.Radius // 50 - 10 = 40
	if ball.X != expectedX {
		t.Errorf("Deep penetration position incorrect. Expected %f, got %f", expectedX, ball.X)
	}
}
