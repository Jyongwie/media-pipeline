import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RenderJobResponse, RenderService } from './services/render';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <div class="dashboard-container">
      
      <header class="dashboard-header">
        <div>
          <h2>Media Pipeline Control</h2>
          <p class="subtitle">Distributed Rendering Architecture</p>
        </div>
        <div class="status-indicator">
          <span class="pulse"></span> API Online
        </div>
      </header>
      
      <div class="card">
        <h3>Queue New Asset</h3>
        <form [formGroup]="jobForm" (ngSubmit)="onSubmit()" class="form-grid">
          
          <div class="form-group">
            <label>Asset Type</label>
            <select formControlName="asset_type">
              <option value="photorealistic-video">Photorealistic Video</option>
              <option value="tactile-audio-track">Tactile Audio Track</option>
            </select>
          </div>

          <div class="form-group">
            <label>Resolution</label>
            <select formControlName="resolution">
              <option value="1080p">1080p Standard</option>
              <option value="4K">4K Ultra HD</option>
            </select>
          </div>

          <div class="form-group full-width">
            <label>Lighting Profile</label>
            <input type="text" formControlName="lighting_profile" placeholder="e.g., warm natural lighting">
          </div>

          <div class="form-group">
            <label>Camera Effect</label>
            <input type="text" formControlName="camera_effect" placeholder="e.g., cinematic deep bokeh">
          </div>

          <div class="form-group">
            <label>Audio Sensitivity</label>
            <input type="text" formControlName="audio_sensitivity" placeholder="e.g., fast tapping focus">
          </div>

          <div class="form-group full-width" style="margin-top: 10px;">
            <button type="submit" [disabled]="!jobForm.valid" class="submit-btn">
              Initialize Render Sequence
            </button>
          </div>
        </form>
      </div>

      <div *ngIf="successMessage" class="toast">
        <svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none">
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
        {{ successMessage }}
      </div>

      <div class="card" style="margin-top: 30px;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
          <h3>Active Render Queue</h3>
          <button (click)="loadJobs()" class="refresh-btn">⟳ Refresh</button>
        </div>

        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Asset Type</th>
                <th>Resolution</th>
                <th>Status</th>
                <th>Submitted</th>
              </tr>
            </thead>
            <tbody>
              <tr *ngFor="let job of jobs">
                <td>{{ job.AssetType }}</td>
                <td><span class="badge">{{ job.Resolution }}</span></td>
                <td>
                  <span class="status-badge" [class.queued]="job.Status === 'queued'">
                    {{ job.Status | uppercase }}
                  </span>
                </td>
                <td style="color: #64748b; font-size: 12px;">
                  {{ job.CreatedAt | date:'shortTime' }}
                </td>
              </tr>
              <tr *ngIf="jobs.length === 0">
                <td colspan="4" style="text-align: center; color: #64748b;">No jobs in queue.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  `,
  styles: [
    `:host {
      display: block;
      min-height: 100vh;
      background-color: #0b0f19; /* Deep dark background */
      color: #e2e8f0;
      font-family: 'Inter', system-ui, sans-serif;
      padding: 40px 20px;
    }

    .dashboard-container {
      max-width: 700px;
      margin: 0 auto;
    }

    .dashboard-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 30px;
      border-bottom: 1px solid #1e293b;
      padding-bottom: 20px;
    }

    h2 { margin: 0; font-size: 24px; font-weight: 600; color: #f8fafc; }
    h3 { margin: 0 0 20px 0; font-size: 18px; font-weight: 500; border-bottom: 1px solid #1e293b; padding-bottom: 10px;}
    .subtitle { margin: 5px 0 0 0; font-size: 14px; color: #64748b; }

    .status-indicator {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 12px;
      color: #10b981;
      background: rgba(16, 185, 129, 0.1);
      padding: 6px 12px;
      border-radius: 20px;
      border: 1px solid rgba(16, 185, 129, 0.2);
    }

    .pulse {
      width: 8px; height: 8px;
      background-color: #10b981;
      border-radius: 50%;
      box-shadow: 0 0 8px #10b981;
      animation: pulse-animation 2s infinite;
    }

    @keyframes pulse-animation {
      0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); }
      70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(16, 185, 129, 0); }
      100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
    }

    .card {
      background-color: #111827;
      border: 1px solid #1e293b;
      border-radius: 12px;
      padding: 25px;
      box-shadow: 0 10px 25px rgba(0,0,0,0.5);
    }

    .form-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 20px;
    }

    .full-width {
      grid-column: 1 / -1;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    label {
      font-size: 12px;
      font-weight: 600;
      color: #94a3b8;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    input, select {
      background-color: #0f172a;
      border: 1px solid #334155;
      color: #f8fafc;
      padding: 12px;
      border-radius: 6px;
      font-size: 14px;
      transition: all 0.2s;
    }

    input:focus, select:focus {
      outline: none;
      border-color: #3b82f6;
      box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
    }

    .submit-btn {
      background: linear-gradient(135deg, #2563eb, #1d4ed8);
      color: white;
      border: none;
      padding: 14px;
      border-radius: 6px;
      font-weight: 600;
      font-size: 14px;
      cursor: pointer;
      transition: transform 0.1s, box-shadow 0.2s;
    }

    .submit-btn:hover:not([disabled]) {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
    }

    .submit-btn:disabled {
      background: #1e293b;
      color: #475569;
      cursor: not-allowed;
    }

    .toast {
      margin-top: 25px;
      padding: 16px;
      background: rgba(16, 185, 129, 0.1);
      border: 1px solid rgba(16, 185, 129, 0.2);
      color: #34d399;
      border-radius: 8px;
      display: flex;
      align-items: center;
      gap: 12px;
      font-weight: 500;
      animation: slide-up 0.3s ease-out;
    }

    @keyframes slide-up {
      from { opacity: 0; transform: translateY(10px); }
      to { opacity: 1; transform: translateY(0); }
    }
    .table-container { overflow-x: auto; }
    table { width: 100%; border-collapse: collapse; text-align: left; }
    th { color: #94a3b8; font-size: 12px; text-transform: uppercase; padding: 12px; border-bottom: 1px solid #1e293b; }
    td { padding: 14px 12px; border-bottom: 1px solid #1e293b; font-size: 14px; }
    
    .badge { background: #1e293b; padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: 600; }
    
    .status-badge { padding: 4px 8px; border-radius: 20px; font-size: 11px; font-weight: 700; }
    .status-badge.queued { background: rgba(245, 158, 11, 0.1); color: #fbbf24; border: 1px solid rgba(245, 158, 11, 0.2); }
    
    .refresh-btn { background: transparent; border: 1px solid #334155; color: #94a3b8; padding: 6px 12px; border-radius: 4px; cursor: pointer; transition: 0.2s; }
    .refresh-btn:hover { background: #1e293b; color: #f8fafc; }`
  ]
})
export class AppComponent implements OnInit {
  private fb = inject(FormBuilder);
  private renderService = inject(RenderService);

  successMessage = '';
  jobs: RenderJobResponse[] = []; // array to hold database record

  // We define our strict form model here
  jobForm = this.fb.group({
    asset_type: ['photorealistic-video', Validators.required],
    resolution: ['4K', Validators.required],
    lighting_profile: ['warm natural lighting', Validators.required],
    camera_effect: ['cinematic deep bokeh', Validators.required],
    audio_sensitivity: ['high-sensitivity-tactile', Validators.required]
  });
  // fetch job on initial load
  ngOnInit() {
    this.loadJobs();
  }

  loadJobs() {
    this.renderService.getJobs().subscribe({
      next: (data) => this.jobs = data || [],
      error: (err) => console.error('Failed to load jobs', err)
    });
  }

  onSubmit() {
    if (this.jobForm.valid) {
      this.renderService.submitJob(this.jobForm.value as any).subscribe({
        next: (response) => {
          this.successMessage = 'Success! ' + response.status;
          this.jobForm.reset(); // Clear the form after success
          this.loadJobs(); // automatically refresh the table when a new job is added

          setTimeout(() => this.successMessage = "", 4000)
        },
        error: (err) => {
          console.error('Submission failed', err);
          alert('Failed to queue job. Is the Go backend running?');
        }
      });
    }
  }
}