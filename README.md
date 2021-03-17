# goda

A simple synth music app using <https://github.com/fr3fou/gusic>

- [x] Render multiple keys
- [ ] Support for pressing several keys at once / play more than 1 note at once
  - [ ] Remove / shift the samples from the buffer as they get played
        - Use a channel instead of an array buffer?
  - [ ] Add up the "conflicting" (the ones that match) samples when copying the new note into the buffer
        - Would have to change from using the `copy` function to a manual process
- [x] Configurable generators
- [ ] Configurable ADSR
- [ ] Note duration based on hold duration
- [x] Better piano keys
- [ ] MIDI Input
- [ ] Octave labels
- [ ] Octave changer
