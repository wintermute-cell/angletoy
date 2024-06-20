package audio

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type playlist []string

type audio struct {
	// volume
	is_mute       bool
	global_volume float32
	music_volume  float32
	sfx_volume    float32

	// tracks
	music_tracks    map[string]rl.Music
	music_playlists map[string]playlist
	sfx_tracks      map[string]rl.Sound

	// playback
	music_fade_secs         float32
	play_timer              float32 // we need this, since GetMusicTimePlayed is a little inaccurate and resets to 0 too early
	curr_playing_music_name string
	curr_playlist           string

	is_music_waiting   bool
	waiting_music_name string // this stores the name of a track that is scheduled to play next

	is_delaying bool    // flag indicating if we're in the delay state between tracks
	delay_timer float32 // a timer to count the delay
	delay_secs  float32 // configurable delay time in seconds between fade-out and fade-in

	force_fade_out     bool
	force_fade_out_len float32 // this is a helper, replacing GetMusicTimeLength with a shorter time.
}

var a audio

func InitAudio() {
	rl.InitAudioDevice()
	if !rl.IsAudioDeviceReady() {
		logging.Fatal("Failed to InitAudioDevice!")
	}
	a = audio{
		is_mute:            false,
		global_volume:      1.0,
		music_volume:       1.0,
		sfx_volume:         1.0,
		music_tracks:       make(map[string]rl.Music),
		music_playlists:    make(map[string]playlist),
		sfx_tracks:         make(map[string]rl.Sound),
		music_fade_secs:    5.0,
		play_timer:         0.0,
		is_music_waiting:   false,
		waiting_music_name: "",
		is_delaying:        false,
		delay_secs:         0.0,
		force_fade_out:     false,
		force_fade_out_len: 0.0,
	}
}

func DeinitAudio() {
	// unload all audio tracks from memory
	for _, sound := range a.sfx_tracks {
		rl.UnloadSound(sound)
	}
	for _, music := range a.music_tracks {
		rl.UnloadMusicStream(music)
	}

	rl.CloseAudioDevice()
}

func Update() {
	if !a.is_mute {

		// Handle the delay between tracks
		if a.is_delaying {
			a.delay_timer += rl.GetFrameTime()
			if a.delay_timer >= a.delay_secs {
				a.is_delaying = false
				a.delay_timer = 0.0
				a.play_timer = 0.0

				// Start playing and fading in the next track
				new_mus := getNextMusicTrack()
				a.curr_playing_music_name = new_mus
				rl.PlayMusicStream(a.music_tracks[new_mus])
				logging.Info("Starting to play music: %v", new_mus)
			}
			return // No more processing needed if we're delaying
		}

		// if were currently playing music...
		if curr_mus, ok := a.music_tracks[a.curr_playing_music_name]; ok {
			rl.UpdateMusicStream(curr_mus)
			a.play_timer += rl.GetFrameTime()

			// get the remaining playtime
			time_remaining := rl.GetMusicTimeLength(curr_mus) - a.play_timer

			// fade in the current music
			if a.play_timer <= a.music_fade_secs {
				rl.SetMusicVolume(curr_mus, a.music_volume*a.global_volume*(a.play_timer/a.music_fade_secs))
			} else if time_remaining <= a.music_fade_secs || a.force_fade_out {
				// handle the fading out logic here...
				fade_time_remaining := rl.GetMusicTimeLength(curr_mus) - a.play_timer
				if a.force_fade_out { // use an alternative calculation if forced to fade out
					fade_time_remaining = a.force_fade_out_len - a.play_timer
				}
				rl.SetMusicVolume(curr_mus, a.music_volume*a.global_volume*util.Max((fade_time_remaining/a.music_fade_secs), 0.01))

				// when the fade out is complete
				if fade_time_remaining <= 0.01 {
					logging.Info("Finished fading out music: %v", a.curr_playing_music_name)
					rl.StopMusicStream(curr_mus)
					a.force_fade_out = false
					a.play_timer = 0.0
					a.is_delaying = true
				}
			} else {
				// if neither fading in, nor out apply the configured music
				// volume (in case configuration changes)
				rl.SetMusicVolume(curr_mus, a.music_volume*a.global_volume)
			}
		} else if _, ok := a.music_playlists[a.curr_playlist]; ok {
			// if we have a playlist, but are not playing
			a.curr_playing_music_name = getNextMusicTrack()
			rl.PlayMusicStream(a.music_tracks[a.curr_playing_music_name])
			logging.Info("Dry starting playlist: %v", a.curr_playlist)
		}

	}
}

func getNextMusicTrack() string {
	// if there is a track waiting, prioritize that
	if a.is_music_waiting {
		a.is_music_waiting = false
		logging.Info("Next music from waiting: %v ", a.waiting_music_name)
		return a.waiting_music_name
	}

	// select random name from a.curr_playlist
	if p, ok := a.music_playlists[a.curr_playlist]; ok {
		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		m := p[rand.Intn(len(p))]
		if m == a.curr_playing_music_name {
			// try once more to somewhat avoid playing the same track
			m = p[rand.Intn(len(p))]
		}
		logging.Info("Next music: %v from playlist: %v", m, a.curr_playlist)
		return m
	} else {
		logging.Error("Current playlist invalid, attempted to select from invalid playlist: %v", a.curr_playlist)
		return ""
	}
}

func forceFadeOutNow() {
	a.force_fade_out = true
	a.force_fade_out_len = rl.GetMusicTimePlayed(a.music_tracks[a.curr_playing_music_name]) + a.music_fade_secs
}

func setWaitingMusic(name string) {
	a.is_music_waiting = true
	a.waiting_music_name = name
}

// ----------------
//       API      |
// ----------------

// LOADING TRACKS & PLAYLISTS

// Load a music file from the given path, register it with the given name.
func RegisterMusic(name, path string) {
	if _, ok := a.music_tracks[name]; ok {
		logging.Warning("Tried to register music track for a name that already exists: %v", name)
		return
	}
	m := rl.LoadMusicStream(path)
	if m.Stream == (rl.AudioStream{}) {
		logging.Error("Failed to load music audio stream for name: \"%v\" and path: %v", name, path)
	}
	m.Looping = false // disable looping by default, this causes issues with fading tracks and looping is done by our player anyway.
	a.music_tracks[name] = m
}

// Load a sound file from the given path, register it with the given name.
func RegisterSound(name, path string) {
	if _, ok := a.sfx_tracks[name]; ok {
		logging.Warning("Tried to register sfx track for a name that already exists: %v", name)
		return
	}
	s := rl.LoadSound(path)
	if s.Stream == (rl.AudioStream{}) {
		logging.Error("Failed to load sound audio stream for name: \"%v\" and path: %v", name, path)
	}
	a.sfx_tracks[name] = s
}

func CreatePlaylist(name string, p []string) {
	if _, ok := a.music_playlists[name]; ok {
		logging.Warning("Tried to register playlist for a name that already exists: %v", name)
		return
	}
	a.music_playlists[name] = p
}

// CONFIGURATION

// Mute the Audio
func Mute() {
	a.is_mute = true
}

// Unmute the Audio
func Unmute() {
	a.is_mute = true
}

// Toggle the Mute State
func ToggleMute() {
	a.is_mute = !a.is_mute
}

// Set the Global Volume to a value between 0.0 and 1.0
func SetGlobalVolume(new_volume float32) {
	new_volume = util.Clamp(new_volume, 0.0, 1.0)
	a.global_volume = new_volume
}

// Get the Global Volume
func GetGlobalVolume() float32 {
	return a.global_volume
}

// Set the Music Volume to a value between 0.0 and 1.0
func SetMusicVolume(new_volume float32) {
	new_volume = util.Clamp(new_volume, 0.0, 1.0)
	a.music_volume = new_volume
}

// Get the Music Volume
func GetMusicVolume() float32 {
	return a.music_volume
}

// Set the SFX Volume to a value between 0.0 and 1.0
func SetSFXVolume(new_volume float32) {
	new_volume = util.Clamp(new_volume, 0.0, 1.0)
	a.sfx_volume = new_volume
}

// Get the SFX Volume
func GetSFXVolume() float32 {
	return a.sfx_volume
}

// Set the Fade Time for Music Tracks in seconds
func SetMusicFade(fade_secs float32) {
	a.music_fade_secs = fade_secs
}

// Get the Fade Time for Music Tracks in seconds
func GetMusicFade() float32 {
	return a.music_fade_secs
}

// PLAYBACK

// Play a sound that has been registered with "name"
func PlaySound(name string) {
	PlaySoundEx(name, 1.0, 1.0, 0.5)
}

/*
PlaySound with extended parameters.

- name: the name of the sound to play

- volume: the volume to play the sound at (0.0 - 1.0)

- pitch: the pitch to play the sound at (default: 1.0)

- pan: the pan to play the sound at (default: 0.5)
*/
func PlaySoundEx(name string, volume, pitch, pan float32) {
	if s, ok := a.sfx_tracks[name]; ok {
		rl.SetSoundPitch(s, pitch)
		rl.SetSoundPan(s, pan)
		rl.SetSoundVolume(s, a.sfx_volume*a.global_volume*volume)

		rl.PlaySound(s)

		// TODO: how do i handle this
		//// reset sound properties
		//rl.SetSoundVolume(s, a.sfx_volume*a.global_volume)
		//rl.SetSoundPitch(s, 1.0)
		//rl.SetSoundPan(s, 0.5)
	} else {
		logging.Warning("Attempted to play sound that is not registered: %v", name)
	}
}

// Instantly start playing a music track that has been registered with "name".
func PlayMusicNow(name string) {
	if _, ok := a.music_tracks[name]; !ok {
		logging.Warning("Attempted to play music that is not registered: %v", name)
		return
	}
	// stop the currently playing track
	rl.StopMusicStream(a.music_tracks[a.curr_playing_music_name])
	// set the correct volume
	rl.SetMusicVolume(a.music_tracks[name], a.music_volume*a.global_volume)
	a.curr_playing_music_name = name
	rl.PlayMusicStream(a.music_tracks[name])
	logging.Info("Playing music now: %v", name)
}

// Start playing a music track that has been registered with "name" after
// fading out the currently playing track.
func PlayMusicNowFade(name string) {
	if _, ok := a.music_tracks[name]; !ok {
		logging.Warning("Attempted to play music that is not registered: %v", name)
		return
	}
	setWaitingMusic(name)
	forceFadeOutNow()
	logging.Info("Force fading into music now: %v", name)
}

// Sets the current playlist to the specified name, optionally fading out the
// currently playing track (should there be one playing).
func SetCurrentPlaylist(name string, fade_current bool) {
	a.curr_playlist = name

	// if music is playing, fade it out
	if _, ok := a.music_tracks[a.curr_playing_music_name]; ok && fade_current {
		forceFadeOutNow()
	}
	logging.Info("Current playlist set to: %v", name)
}
