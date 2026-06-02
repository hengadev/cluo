# ADR-0005: S3 Storage Path Segregation by MediaFile Type

## Status
Accepted

## Context
A Case accumulates three categories of MediaFile during an investigation: audio recordings, photos, and videos. All three are stored in the same S3 bucket. Without explicit path segregation, all MediaFiles for a case would land under the same prefix regardless of type, making it harder to manage storage, apply per-type lifecycle policies, or audit what was captured.

## Decision
MediaFiles are stored under type-segregated S3 key prefixes within the same bucket. Audio files (Recordings) live under a dedicated prefix, separate from images and videos.

The term **Recording** is reserved for audio MediaFiles (`type: "audio"`). It is a named subtype of MediaFile and the primary capture target of `cluo_mobile`.

## Consequences
- Storage for Recordings can be managed, monitored, and lifecycle-governed independently from photo/video storage.
- Any new upload path must respect the segregation — uploading an audio file to the image prefix (or vice versa) is a bug.
- The mobile app treats Recordings as a distinct concept and does not mix them with image/video upload flows.
