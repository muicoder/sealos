/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DebtStatusNormal DebtStatusType = "Normal"
	DebtStatusSmall  DebtStatusType = "Small"
	DebtStatusMedium DebtStatusType = "Medium"
	DebtStatusLarge  DebtStatusType = "Large"
	DebtPrefix                      = "debt-"
	NormalPrice                     = 0
	// SmallBlockWaitSecond 3 days
	SmallBlockWaitSecond = 3 * 24 * 60 * 60
	// MediumBlockWaitSecond 4 days
	MediumBlockWaitSecond = 4 * 24 * 60 * 60
)

type DebtStatusType string

var DefaultDebtConfig = map[string]int64{
	"Normal": NormalPrice,
	"Small":  SmallBlockWaitSecond,
	"Medium": MediumBlockWaitSecond,
}

// DebtSpec defines the desired state of Debt
type DebtSpec struct {
	UserName string `json:"userName,omitempty"`
}

// DebtStatus defines the observed state of Debt
type DebtStatus struct {
	LastUpdateTimestamp int64          `json:"lastUpdateTimestamp,omitempty"`
	AccountDebtStatus   DebtStatusType `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Debt is the Schema for the debts API
type Debt struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DebtSpec   `json:"spec,omitempty"`
	Status DebtStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DebtList contains a list of Debt
type DebtList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Debt `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Debt{}, &DebtList{})
}