import { Injectable, isDevMode } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

// We define an interface that perfectly matches our Go backend struct
export interface RenderJobRequest {
  asset_type: string;
  resolution: string;
  lighting_profile: string;
  camera_effect: string;
  audio_sensitivity: string;
}
// for data coming back from go
export interface RenderJobResponse {
  ID: string;
  AssetType: string;
  Resolution: string;
  LightingProfile: string;
  CameraEffect: string;
  AudioSensitivity: string;
  Status: string;
  CreatedAt: string;
}

@Injectable({
  providedIn: 'root'
})
export class RenderService {
  private apiUrl = isDevMode()? '/api/jobs':'https://media-backend-9v1v.onrender.com/api/jobs';
  constructor(private http: HttpClient) { }

  submitJob(job: RenderJobRequest): Observable<{status: string}> {
    // Notice we just call /api/jobs. The proxy handles the rest!
    return this.http.post<{status: string}>(this.apiUrl, job);
  }
  // NEW: Fetch all jobs
  getJobs(): Observable<RenderJobResponse[]> {
    return this.http.get<RenderJobResponse[]>(this.apiUrl);
  }
}