package model

import "testing"

func TestDeriveReleaseStatusFromSteps(t *testing.T) {
	tests := []struct {
		name          string
		releaseAction string
		currentStatus ReleaseStatus
		steps         []ReleaseStep
		want          ReleaseStatus
	}{
		{
			name:          "rollback succeeded becomes rolled back",
			releaseAction: "Rollback",
			currentStatus: ReleaseRunning,
			steps: []ReleaseStep{
				{Name: "apply rollback manifests", Status: StepSucceeded},
			},
			want: ReleaseRolledBack,
		},
		{
			name:          "failed step becomes failed",
			releaseAction: "Upgrade",
			currentStatus: ReleaseRunning,
			steps: []ReleaseStep{
				{Name: "apply manifests", Status: StepFailed},
			},
			want: ReleaseFailed,
		},
		{
			name:          "running step becomes running",
			releaseAction: "Install",
			currentStatus: ReleasePending,
			steps: []ReleaseStep{
				{Name: "apply install manifests", Status: StepRunning},
			},
			want: ReleaseRunning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeriveReleaseStatusFromSteps(tt.releaseAction, tt.currentStatus, tt.steps); got != tt.want {
				t.Fatalf("DeriveReleaseStatusFromSteps() = %q, want %q", got, tt.want)
			}
		})
	}
}
