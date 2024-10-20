graph TD
A[Screen Capture Module] -->|Raw Frames| B[Frame Buffer]
B -->|Raw Frames| C[Video Encoder]
D[Audio Capture Module] -->|Raw Audio| E[Audio Buffer]
E -->|Raw Audio| F[Audio Encoder]
C -->|Encoded Video| G[RTMP Packager]
F -->|Encoded Audio| G
G -->|RTMP Stream| H[Twitch RTMP Ingest Server]
I[Configuration Manager] -->|Settings| A
I -->|Settings| D
I -->|Settings| C
I -->|Settings| F
I -->|Settings| G
J[User Interface] -->|Commands| K[Stream Controller]
K -->|Control| A
K -->|Control| D
K -->|Control| G
L[Twitch API Client] -->|Stream Key/Info| G
M[Error Handler] -->|Logs| N[Logging System]