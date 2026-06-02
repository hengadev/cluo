# ADR-0006: Audio Recordings excluded from the client portal

Audio Recordings are captured in the field by the Investigator and transcribed by AI; their value to the Client is delivered through the Rapport, not as raw files. Despite being a subtype of MediaFile that can be marked `IsPublished`, audio must never appear in `cluo_web` — neither returned by the portal media endpoint (`GET /token/{token}/media`) nor rendered in the UI.

The alternative — letting the PI choose to publish audio like images or videos — was rejected on two grounds: (1) raw recordings have no interpretive value to a non-investigator and would only create confusion; (2) audio recordings captured during surveillance carry legal sensitivity that makes uncontrolled client delivery unacceptable.

The backend must filter out `type = "audio"` from the portal media response, and the frontend must not render audio even if a future API change were to include it.
