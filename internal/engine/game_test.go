package engine

import "testing"

func TestAvatarNextAndReturn(t *testing.T) {
	game := newgame()
	g := &game //required methods are defined on the pointer, not the value
	//avatar pool length at this point is 2

	//remove one
	avt := g.nextAvatar()
	if len(g.avatarPool) != 1 {
		t.Errorf("nextAvatar does not remove from pool, expected pool length %d, got %d", 1, len(g.avatarPool))
	}

	//add the avatar back
	g.returnAvatar(avt)
	if len(g.avatarPool) != 2 {
		t.Errorf("returnAvatar does not return to pool, expected pool length %d, got %d", 2, len(g.avatarPool))
	}
}
