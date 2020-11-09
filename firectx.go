package firectx

import (
	"context"

	"cloud.google.com/go/firestore"
)

type clientKeyType struct{}

var clientKey clientKeyType

// WithFirestoreClient must be called before using Firestore client.
func WithFirestoreClient(ctx context.Context, c *firestore.Client) context.Context {
	return context.WithValue(ctx, clientKey, c)
}

func FirestoreClient(ctx context.Context) *firestore.Client {
	c, _ := ctx.Value(clientKey).(*firestore.Client)
	return c
}

type txKeyType struct{}

var txKey txKeyType

func RunTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	return FirestoreClient(ctx).RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return f(txCtx)
	})
}

func transaction(ctx context.Context) (*firestore.Transaction, bool) {
	tx, ok := ctx.Value(txKey).(*firestore.Transaction)
	return tx, ok
}

func Create(ctx context.Context, dr *firestore.DocumentRef, data interface{}) error {
	tx, ok := transaction(ctx)
	if ok {
		return tx.Create(dr, data)
	}

	_, err := dr.Create(ctx, data)
	return err
}

func Delete(ctx context.Context, dr *firestore.DocumentRef, opts ...firestore.Precondition) error {
	tx, ok := transaction(ctx)
	if ok {
		return tx.Delete(dr, opts...)
	}

	_, err := dr.Delete(ctx, opts...)
	return err
}

func Get(ctx context.Context, dr *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	tx, ok := transaction(ctx)
	if ok {
		return tx.Get(dr)
	}

	return dr.Get(ctx)
}

func Set(ctx context.Context, dr *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error {
	tx, ok := transaction(ctx)
	if ok {
		return tx.Set(dr, data, opts...)
	}

	_, err := dr.Set(ctx, data, opts...)
	return err
}

func Update(ctx context.Context, dr *firestore.DocumentRef, data []firestore.Update, opts ...firestore.Precondition) error {
	tx, ok := transaction(ctx)
	if ok {
		return tx.Update(dr, data, opts...)
	}

	_, err := dr.Update(ctx, data, opts...)
	return err
}

func Documents(ctx context.Context, q firestore.Query) *firestore.DocumentIterator {
	tx, ok := transaction(ctx)
	if ok {
		return tx.Documents(q)
	}

	return q.Documents(ctx)
}

func DocumentRefs(ctx context.Context, cr *firestore.CollectionRef) *firestore.DocumentRefIterator {
	tx, ok := transaction(ctx)
	if ok {
		return tx.DocumentRefs(cr)
	}

	return cr.DocumentRefs(ctx)
}

func Collection(ctx context.Context, path string) *firestore.CollectionRef {
	return FirestoreClient(ctx).Collection(path)
}
