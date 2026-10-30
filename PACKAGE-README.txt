HELIOS DRIFT
============

Space trading game. Two modes: Classic (the 1984 rules) and Advanced
(campaign, contracts, shipyard, mining). Every launch gets its own world;
the seed is printed on the title screen and can be replayed with -seed.

RUN
---
Linux:    chmod +x heliosdrift && ./heliosdrift
Windows:  heliosdrift.exe  (heliosdrift-console.exe prints diagnostics)

The game opens a 960x672 window, resizable. F11 toggles fullscreen.

CONTROLS (flight)
-----------------
Arrows / WASD / mouse   pitch and roll
Space / left mouse      lasers
X / Z                   accelerate / brake
M                       missile
E                       ECM
Tab                     energy bomb
L                       docking computer
G                       galactic hyperspace
J                       jump
K                       jettison cargo
H                       hyperspace
F10                     mute

In the station, digits pick menu entries, arrows move, Enter selects,
Escape backs out. A gamepad works everywhere: sticks fly, A/B/D-pad drive
the menus.

Saves and settings live in your home directory:
  .heliosdrift-save1.json ... -save3.json
  .heliosdrift-config.json
  .heliosdrift-keys.json

Takes no external assets: every ship, station and sound is generated at
start-up.
