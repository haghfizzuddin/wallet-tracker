package progress

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/briandowns/spinner"
	"github.com/schollz/progressbar/v3"
)

// Spinner provides a simple spinner for indeterminate progress
type Spinner struct {
	spinner *spinner.Spinner
	message string
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	return &Spinner{
		spinner: s,
		message: message,
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.spinner.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.spinner.Stop()
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.spinner.Suffix = " " + message
}

// Bar provides a progress bar for determinate progress
type Bar struct {
	bar   *progressbar.ProgressBar
	total int
	mu    sync.Mutex
}

// NewBar creates a new progress bar
func NewBar(total int, description string) *Bar {
	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)
	
	return &Bar{
		bar:   bar,
		total: total,
	}
}

// Increment increments the progress bar
func (b *Bar) Increment() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bar.Add(1)
}

// IncrementBy increments the progress bar by a specific amount
func (b *Bar) IncrementBy(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bar.Add(n)
}

// Finish completes the progress bar
func (b *Bar) Finish() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bar.Finish()
}

// UpdateDescription updates the progress bar description
func (b *Bar) UpdateDescription(description string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bar.Describe(description)
}

// MultiProgress manages multiple progress indicators
type MultiProgress struct {
	spinners map[string]*Spinner
	bars     map[string]*Bar
	mu       sync.RWMutex
}

// NewMultiProgress creates a new multi-progress manager
func NewMultiProgress() *MultiProgress {
	return &MultiProgress{
		spinners: make(map[string]*Spinner),
		bars:     make(map[string]*Bar),
	}
}

// AddSpinner adds a new spinner
func (mp *MultiProgress) AddSpinner(id, message string) *Spinner {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	spinner := NewSpinner(message)
	mp.spinners[id] = spinner
	return spinner
}

// AddBar adds a new progress bar
func (mp *MultiProgress) AddBar(id string, total int, description string) *Bar {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	bar := NewBar(total, description)
	mp.bars[id] = bar
	return bar
}

// GetSpinner gets a spinner by ID
func (mp *MultiProgress) GetSpinner(id string) *Spinner {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.spinners[id]
}

// GetBar gets a progress bar by ID
func (mp *MultiProgress) GetBar(id string) *Bar {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.bars[id]
}

// StopAll stops all progress indicators
func (mp *MultiProgress) StopAll() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	for _, spinner := range mp.spinners {
		spinner.Stop()
	}
	
	for _, bar := range mp.bars {
		bar.Finish()
	}
}
