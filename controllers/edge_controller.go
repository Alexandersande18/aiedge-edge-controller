/*
Copyright 2022.

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

package controllers

import (
	"context"
	glog "log"
	aiedgendsllabcnv1 "aiedge-edge-controller/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	AiedgeEdgeTagName = "aiedge/edge"
	KEEdgeTagName         = "node-role.kubernetes.io/edge"
)

// EdgeReconciler reconciles a Edge object
type EdgeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=aiedge.ndsl-lab.cn,resources=edges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aiedge.ndsl-lab.cn,resources=edges/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aiedge.ndsl-lab.cn,resources=edges/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Edge object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *EdgeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var edge aiedgendsllabcnv1.Edge
	if err := r.Get(ctx, req.NamespacedName, &edge); err != nil {
		log.Error(err, "Failed to GET edges")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//DELETION
	myFinalizerName := "aiedge.ndsl-lab.cn/finalizer"
    // examine DeletionTimestamp to determine if object is under deletion
    if edge.ObjectMeta.DeletionTimestamp.IsZero() {
		glog.Println("Not being deleted")
        if !controllerutil.ContainsFinalizer(&edge, myFinalizerName) {
			glog.Println("Adding finalizer")
            controllerutil.AddFinalizer(&edge, myFinalizerName)
            if err := r.Update(ctx, &edge); err != nil {
                return ctrl.Result{}, err
            }
        }
    } else {
        // The object is being deleted
		glog.Println("Being deleted")
        if controllerutil.ContainsFinalizer(&edge, myFinalizerName) {
            if err := r.deleteExternalResources(&edge); err != nil {
                return ctrl.Result{}, err
            }
            controllerutil.RemoveFinalizer(&edge, myFinalizerName)
            if err := r.Update(ctx, &edge); err != nil {
                return ctrl.Result{}, err
            }
        }
        return ctrl.Result{}, nil
    }

	// ADD & UPDATE
	r.updateStatus(ctx, &edge)
	edgeName := edge.GetName()
	r.applyTest(ctx, &edge)
	for _, nodeT := range edge.Spec.Nodes {
		var kNode corev1.Node
		if err := r.Get(ctx, ktypes.NamespacedName{Namespace: "", Name: nodeT}, &kNode); err != nil {
			log.Error(err, "Failed to GET Node" + nodeT)
			return ctrl.Result{}, client.IgnoreAlreadyExists(err)
		} 
		glog.Println(edgeName +  "/" + kNode.Name)
		lb := kNode.GetLabels()
		if _, ok := lb[KEEdgeTagName]; !ok {
			continue
		}
		if label, ok := lb[AiedgeEdgeTagName]; ok && label == edgeName {
			continue
		}
		kNode.Labels[AiedgeEdgeTagName] = edgeName	
		if err := r.Patch(ctx, &kNode, client.Merge); err != nil {
			log.Error(err, "Failed to TAG Node" + nodeT)
		}
	}
	
	var nodeList corev1.NodeList
	// labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"version":""}}
	labelSelector, _ := labels.Parse("node-role.kubernetes.io/edge")
	if err := r.List(ctx, &nodeList, &client.ListOptions{LabelSelector: labelSelector}); err != nil {
		log.Error(err, "Failed to LIST nodes")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}


func (r *EdgeReconciler) updateStatus(ctx context.Context, edge *aiedgendsllabcnv1.Edge ) error {
	log := log.FromContext(ctx)
	edge.Status.EdgeSize = int32(len(edge.Spec.Nodes))
    if err := r.Status().Update(context.Background(), edge); err != nil {
		log.Error(err, "Failed to UPDATE status")
		return err
	}
	return nil
}

func (r *EdgeReconciler) deleteExternalResources(edge *aiedgendsllabcnv1.Edge) error {
    //
    // delete any external resources associated with the cronJob
    //
    // Ensure that delete implementation is idempotent and safe to invoke
    // multiple times for same object.
	return nil
}

func (r *EdgeReconciler) applyTest(ctx context.Context, edge *aiedgendsllabcnv1.Edge) error {
	// log := log.FromContext(ctx)
	// clientObj, err := deserializeFromFile("/home/zqd/k8s-test/aiedge-edge-controller/yamls/test.yaml")
	// if err != nil {
	// 	log.Error(err, "deserializeFromFile Error")
	// 	return err
	// }
	// r.Create(ctx, clientObj)
	return nil
}


// SetupWithManager sets up the controller with the Manager.
func (r *EdgeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiedgendsllabcnv1.Edge{}).
		Complete(r)
}
