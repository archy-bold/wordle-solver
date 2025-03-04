package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	kb                = newKeyboard()
	validWords        = []string{"group", "prank", "spare", "tapir"}
	validWords2       = []string{"at", "ta"}
	gameTapirStart    = &game{false, 0, "tapir", make(Grid, 6), &validWords, kb, 200}
	gameTapirFinished = &game{true, 4, "tapir", Grid{
		{GridCell{"g", STATUS_WRONG}, GridCell{"r", STATUS_INCORRECT}, GridCell{"o", STATUS_WRONG}, GridCell{"u", STATUS_WRONG}, GridCell{"p", STATUS_INCORRECT}},
		{GridCell{"p", STATUS_INCORRECT}, GridCell{"r", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}, GridCell{"n", STATUS_WRONG}, GridCell{"k", STATUS_WRONG}},
		{GridCell{"s", STATUS_WRONG}, GridCell{"p", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}, GridCell{"r", STATUS_INCORRECT}, GridCell{"e", STATUS_WRONG}},
		{GridCell{"t", STATUS_CORRECT}, GridCell{"a", STATUS_CORRECT}, GridCell{"p", STATUS_CORRECT}, GridCell{"i", STATUS_CORRECT}, GridCell{"r", STATUS_CORRECT}},
		nil,
		nil,
	}, &validWords, kb, 200}
	gameAtStart    = &game{false, 0, "at", make(Grid, 1), &validWords2, kb, 200}
	gameAtFinished = &game{false, 1, "at", Grid{
		{GridCell{"t", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}},
	}, &validWords2, kb, 200}
)

var createGameTests = map[string]struct {
	answer     string
	tries      int
	validWords *[]string
	gameNum    int
	expected   *game
}{
	"5 letter, 6 tries":   {"tapir", 6, &validWords, 200, &game{false, 0, "tapir", Grid{nil, nil, nil, nil, nil, nil}, &validWords, kb, 200}},
	"3 letter, 3 tries":   {"bat", 3, &validWords2, 199, &game{false, 0, "bat", Grid{nil, nil, nil}, &validWords2, kb, 199}},
	"5 letter, uppercase": {"TAPIR", 6, &validWords, 200, &game{false, 0, "tapir", Grid{nil, nil, nil, nil, nil, nil}, &validWords, kb, 200}},
}

func Test_CreateGame(t *testing.T) {
	for tn, tt := range createGameTests {
		g := CreateGame(tt.answer, tt.tries, tt.validWords, tt.gameNum)

		assert.Equalf(t, tt.expected, g, "Expected game to match for test '%s'", tn)
	}
}

var gamePlayTests = map[string]struct {
	g            *game
	tries        []string
	expected     []bool
	expectedErr  string
	expectedGrid Grid
}{
	"5-letter, won": {
		g:            gameTapirStart,
		tries:        []string{"group", "prank", "spare", "tapir"},
		expected:     []bool{false, false, false, true},
		expectedGrid: gameTapirFinished.grid,
	},
	"5-letter, mixed case": {
		g:            gameTapirStart,
		tries:        []string{"grOUp", "PRAnk", "sPaRE", "TAPIR"},
		expected:     []bool{false, false, false, true},
		expectedGrid: gameTapirFinished.grid,
	},
	"5-letter, try 4-letter word": {
		g:           gameTapirStart,
		tries:       []string{"tape"},
		expectedErr: "The entered word length is wrong, should be: 5",
	},
	"5-letter, try 6-letter word": {
		g:           gameTapirStart,
		tries:       []string{"strong"},
		expectedErr: "The entered word length is wrong, should be: 5",
	},
	"5-letter, try invalid word": {
		g:           gameTapirStart,
		tries:       []string{"scrap"},
		expectedErr: ErrInvalidWord.Error(),
	},
	"2-letter, lost": {
		g:            gameAtStart,
		tries:        []string{"ta"},
		expected:     []bool{false},
		expectedGrid: gameAtFinished.grid,
	},
}

func Test_game_Play(t *testing.T) {
	for tn, tt := range gamePlayTests {
		// Copy first
		g := &game{tt.g.complete, tt.g.attempts, tt.g.answer, make(Grid, len(tt.g.grid)), tt.g.validWords, newKeyboard(), tt.g.gameNum}
		for i, row := range tt.g.grid {
			copy(g.grid[i], row)
		}
		for i, word := range tt.tries {
			res, err := g.Play(word)

			// Make the assertions
			if tt.expectedErr != "" {
				assert.Falsef(t, res, "Expected res false for test '%s', try %d", tn, i)
				assert.Errorf(t, err, "Expected error to match for test '%s', try %d", tn, i)
			} else {
				assert.NoErrorf(t, err, "Expected nil error for test '%s', try %d", tn, i)
				assert.Equalf(t, tt.expected[i], res, "Expected play outcome to match for test '%s', try %d", tn, i)
				assert.Equalf(t, tt.expected[i], g.complete, "Expected complete to match for test '%s', try %d", tn, i)
				assert.Equalf(t, tt.expectedGrid[i], g.grid[i], "Expected grid row to match for test '%s', try %d", tn, i)
				assert.Equalf(t, i+1, g.attempts, "Expected attempts to match for test '%s', try %d", tn, i)
			}
		}
	}
}

var gameHasEndedTests = map[string]struct {
	g        *game
	expected bool
}{
	"5-letter start":    {gameTapirStart, false},
	"5-letter finished": {gameTapirFinished, true},
	"2-letter start":    {gameAtStart, false},
	"2-letter finished": {gameAtFinished, true},
}

func Test_game_HasEnded(t *testing.T) {
	for tn, tt := range gameHasEndedTests {
		assert.Equalf(t, tt.expected, tt.g.HasEnded(), "Expected result to match for test '%s'", tn)
	}
}

var gameGetScoreTests = map[string]struct {
	g             *game
	expectedScore int
	expectedOf    int
}{
	"5-letter start":    {gameTapirStart, 0, 6},
	"5-letter finished": {gameTapirFinished, 4, 6},
	"2-letter start":    {gameAtStart, 0, 1},
	"2-letter finished": {gameAtFinished, 1, 1},
}

func Test_game_GetScore(t *testing.T) {
	for tn, tt := range gameGetScoreTests {
		score, of := tt.g.GetScore()

		assert.Equalf(t, tt.expectedScore, score, "Expected score to match for test '%s'", tn)
		assert.Equalf(t, tt.expectedOf, of, "Expected of to match for test '%s'", tn)
	}
}

var gameGetLastPlayTests = map[string]struct {
	g        *game
	expected []GridCell
}{
	"5-letter start":    {gameTapirStart, nil},
	"5-letter finished": {gameTapirFinished, gameTapirFinished.grid[3]},
	"2-letter start":    {gameAtStart, nil},
	"2-letter finished": {gameAtFinished, gameAtFinished.grid[0]},
}

func Test_game_GetLastPlay(t *testing.T) {
	for tn, tt := range gameGetLastPlayTests {
		assert.Equalf(t, tt.expected, tt.g.GetLastPlay(), "Expected result to match for test '%s'", tn)
	}
}

var gameOutputForConsoleTests = map[string]struct {
	g        *game
	expected string
}{
	"5-letter start": {
		g:        gameTapirStart,
		expected: "\n       -------\n       -------\n" + defaulKBOutput,
	},
	"5-letter finished": {
		g: gameTapirFinished,
		expected: "\n       -------\n" +
			"       |G" + COLOUR_RESET + COLOUR_YELLOW + "R" + COLOUR_RESET + "O" + COLOUR_RESET + "U" + COLOUR_RESET + COLOUR_YELLOW + "P" + COLOUR_RESET + "|\n" +
			"       |" + COLOUR_YELLOW + "P" + COLOUR_RESET + COLOUR_YELLOW + "R" + COLOUR_RESET + COLOUR_YELLOW + "A" + COLOUR_RESET + "N" + COLOUR_RESET + "K" + COLOUR_RESET + "|\n" +
			"       |S" + COLOUR_RESET + COLOUR_YELLOW + "P" + COLOUR_RESET + COLOUR_YELLOW + "A" + COLOUR_RESET + COLOUR_YELLOW + "R" + COLOUR_RESET + "E" + COLOUR_RESET + "|\n" +
			"       |" + COLOUR_GREEN + "T" + COLOUR_RESET + COLOUR_GREEN + "A" + COLOUR_RESET + COLOUR_GREEN + "P" + COLOUR_RESET + COLOUR_GREEN + "I" + COLOUR_RESET + COLOUR_GREEN + "R" + COLOUR_RESET + "|\n" +
			"       -------\n" +
			defaulKBOutput,
	},
	"2-letter start": {
		g:        gameAtStart,
		expected: "\n       ----\n       ----\n" + defaulKBOutput,
	},
	"2-letter finished": {
		g:        gameAtFinished,
		expected: "\n       ----\n       |" + COLOUR_YELLOW + "T" + COLOUR_RESET + COLOUR_YELLOW + "A" + COLOUR_RESET + "|\n       ----\n" + defaulKBOutput,
	},
}

func Test_game_OutputForConsole(t *testing.T) {
	for tn, tt := range gameOutputForConsoleTests {
		assert.Equalf(t, tt.expected, tt.g.OutputForConsole(), "Expected result to match for test '%s'", tn)
	}
}

var gameOutputToShareTests = map[string]struct {
	g        *game
	expected string
}{
	"5-letter start": {
		g:        gameTapirStart,
		expected: "Wordle 200 0/6\n\n\n",
	},
	"5-letter finished": {
		g: gameTapirFinished,
		expected: "Wordle 200 4/6\n\n" +
			"⬜🟨⬜⬜🟨\n" +
			"🟨🟨🟨⬜⬜\n" +
			"⬜🟨🟨🟨⬜\n" +
			"🟩🟩🟩🟩🟩\n\n",
	},
	"2-letter start": {
		g:        gameAtStart,
		expected: "Wordle 200 0/1\n\n\n",
	},
	"2-letter finished": {
		g:        gameAtFinished,
		expected: "Wordle 200 X/1\n\n🟨🟨\n\n",
	},
}

func Test_game_OutputToShare(t *testing.T) {
	for tn, tt := range gameOutputToShareTests {
		assert.Equalf(t, tt.expected, tt.g.OutputToShare(), "Expected result to match for test '%s'", tn)
	}
}
