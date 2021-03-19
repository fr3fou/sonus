# goda

A simple synth music app using <https://github.com/fr3fou/gusic>

- [x] Render multiple keys
- [ ] Support for pressing several keys at once / play more than 1 note at once
  - [ ] Remove / shift the samples from the buffer as they get played - Use a channel instead of an array buffer?
  - [ ] Add up the "conflicting" (the ones that match) samples when copying the new note into the buffer - Would have to change from using the `copy` function to a manual process
- [x] Configurable generators
- [ ] ADSR
  - [ ] Adjust ratios
  - [ ] Handle triggering and releasing (<https://github.com/velipso/adsrnode#triggering-and-releasing>)
- [ ] Note duration based on hold duration
- [ ] Fix clipping
- [ ] Support for continuous notes (no gap between individual notes)
- [ ] Support for singular notes
- [x] Better piano keys
- [ ] Reverb
- [ ] MIDI Input
- [ ] Octave labels
- [ ] Octave changer
