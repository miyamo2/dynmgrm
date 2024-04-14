package tests

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	testTableName = "test_tables"
)

var testData = []map[string]*dynamodb.AttributeValue{
	{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	},
	{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("2"),
		},
		"some_string": {
			S: aws.String("こんにちは"),
		},
		"some_int": {
			N: aws.String("2"),
		},
		"some_float": {
			N: aws.String("2.2"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("GHI"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("こんにちは"),
				},
				{
					N: aws.String("2"),
				},
				{
					N: aws.String("2.2"),
				},
				{
					BOOL: aws.Bool(false),
				},
				{
					B: []byte("GHI"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("こんにちは"),
				},
				"some_number": {
					N: aws.String("2.2"),
				},
				"some_bool": {
					BOOL: aws.Bool(false),
				},
				"some_binary": {
					B: []byte("GHI"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("こんにちは"), aws.String("世界")},
		},
		"some_int_set": {
			NS: []*string{aws.String("2"), aws.String("4")},
		},
		"some_float_set": {
			NS: []*string{aws.String("2.2"), aws.String("4.4")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("GHI"), []byte("JKL")},
		},
		"any": {
			N: aws.String("0"),
		},
	},
	{
		"pk": {
			S: aws.String("Partition2"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	},
	{
		"pk": {
			S: aws.String("Partition2"),
		},
		"sk": {
			N: aws.String("2"),
		},
		"some_string": {
			S: aws.String("こんにちは"),
		},
		"some_int": {
			N: aws.String("2"),
		},
		"some_float": {
			N: aws.String("2.2"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("GHI"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("こんにちは"),
				},
				{
					N: aws.String("2"),
				},
				{
					N: aws.String("2.2"),
				},
				{
					BOOL: aws.Bool(false),
				},
				{
					B: []byte("GHI"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("こんにちは"),
				},
				"some_number": {
					N: aws.String("2.2"),
				},
				"some_bool": {
					BOOL: aws.Bool(false),
				},
				"some_binary": {
					B: []byte("GHI"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("こんにちは"), aws.String("世界")},
		},
		"some_int_set": {
			NS: []*string{aws.String("2"), aws.String("4")},
		},
		"some_float_set": {
			NS: []*string{aws.String("2.2"), aws.String("4.4")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("GHI"), []byte("JKL")},
		},
		"any": {
			N: aws.String("0"),
		},
	},
	{
		"pk": {
			S: aws.String("Partition3"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
	},
	{
		"pk": {
			S: aws.String("Partition3"),
		},
		"sk": {
			N: aws.String("2"),
		},
		"some_string": {
			S: aws.String("こんにちは"),
		},
		"some_int": {
			N: aws.String("2"),
		},
		"some_float": {
			N: aws.String("2.2"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("GHI"),
		},
	},
}

var testDataForTypedList = []map[string]*dynamodb.AttributeValue{
	{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"typed_list": {
			L: []*dynamodb.AttributeValue{
				{
					M: map[string]*dynamodb.AttributeValue{
						"some_string": {
							S: aws.String("Hello"),
						},
						"some_int": {
							N: aws.String("1"),
						},
						"some_float": {
							N: aws.String("1.1"),
						},
						"some_bool": {
							BOOL: aws.Bool(true),
						},
						"some_binary": {
							B: []byte("ABC"),
						},
						"some_list": {
							L: []*dynamodb.AttributeValue{
								{
									S: aws.String("Hello"),
								},
								{
									N: aws.String("1"),
								},
								{
									N: aws.String("1.1"),
								},
								{
									BOOL: aws.Bool(true),
								},
								{
									B: []byte("ABC"),
								},
							},
						},
						"some_map": {
							M: map[string]*dynamodb.AttributeValue{
								"some_string": {
									S: aws.String("Hello"),
								},
								"some_number": {
									N: aws.String("1.1"),
								},
								"some_bool": {
									BOOL: aws.Bool(true),
								},
								"some_binary": {
									B: []byte("ABC"),
								},
							},
						},
						"some_string_set": {
							SS: []*string{aws.String("Hello"), aws.String("World")},
						},
						"some_int_set": {
							NS: []*string{aws.String("1"), aws.String("2")},
						},
						"some_float_set": {
							NS: []*string{aws.String("1.1"), aws.String("2.2")},
						},
						"some_binary_set": {
							BS: [][]byte{[]byte("ABC"), []byte("DEF")},
						},
					},
				},
			},
		},
	},
}
