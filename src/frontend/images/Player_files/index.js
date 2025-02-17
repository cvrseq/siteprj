const { default: WaveSurfer } = require('wavesurfer.js');

duration = document.querySelector('#duration');
current = document.querySelector('#current');
playPause = document.querySelector('#playPause');

var timeCalculator = function (value) {
  second = Math.floor(value % 60);
  minute = Math.floor((value * 60) % 60);

  if (second < 10) {
    second = '0' + second;
  }

  return minute + ':' + second;
};

// <----------------------------------------------------->

wavesurfer = WaveSurfer.create({
  container: '#wave',
  waveColor: '#cdedff',
  progressColor: '#1aafff',
  height: 48,
  scrollParent: false,
});

wavesurfer.load('./Chris Isaak - Wicked Game (online-audio-converter.com).m4a');
