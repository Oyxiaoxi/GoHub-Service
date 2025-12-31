// Package testutil 提供测试辅助工具
package testutil

import (
	"encoding/json"
	"reflect"
	"testing"
)

// AssertEqual 断言两个值相等
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s\n期望: %+v\n实际: %+v", message, expected, actual)
	}
}

// AssertNotEqual 断言两个值不相等
func AssertNotEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("%s\n不应该相等: %+v", message, expected)
	}
}

// AssertNil 断言值为nil
func AssertNil(t *testing.T, actual interface{}, message string) {
	t.Helper()
	if actual != nil && !reflect.ValueOf(actual).IsNil() {
		t.Errorf("%s\n期望: nil\n实际: %+v", message, actual)
	}
}

// AssertNotNil 断言值不为nil
func AssertNotNil(t *testing.T, actual interface{}, message string) {
	t.Helper()
	if actual == nil || reflect.ValueOf(actual).IsNil() {
		t.Errorf("%s\n期望: 非nil\n实际: nil", message)
	}
}

// AssertTrue 断言条件为真
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf("%s\n期望: true\n实际: false", message)
	}
}

// AssertFalse 断言条件为假
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf("%s\n期望: false\n实际: true", message)
	}
}

// AssertNoError 断言没有错误
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s\n不应该有错误，但得到: %v", message, err)
	}
}

// AssertError 断言有错误
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s\n应该有错误，但没有", message)
	}
}

// AssertContains 断言字符串包含子串
func AssertContains(t *testing.T, str, substr, message string) {
	t.Helper()
	if len(str) == 0 || len(substr) == 0 {
		t.Errorf("%s\n字符串或子串为空", message)
		return
	}
	// 简单的包含检查
	found := false
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("%s\n字符串 '%s' 应该包含 '%s'", message, str, substr)
	}
}

// AssertJSONEqual 断言两个JSON字符串相等（忽略格式差异）
func AssertJSONEqual(t *testing.T, expected, actual string, message string) {
	t.Helper()
	var expectedJSON, actualJSON interface{}
	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		t.Errorf("解析期望的JSON失败: %v", err)
		return
	}
	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		t.Errorf("解析实际的JSON失败: %v", err)
		return
	}
	if !reflect.DeepEqual(expectedJSON, actualJSON) {
		t.Errorf("%s\n期望: %s\n实际: %s", message, expected, actual)
	}
}

// AssertPanic 断言函数会panic
func AssertPanic(t *testing.T, f func(), message string) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s\n应该panic但没有", message)
		}
	}()
	f()
}

// AssertGreaterThan 断言a > b
func AssertGreaterThan(t *testing.T, a, b int, message string) {
	t.Helper()
	if a <= b {
		t.Errorf("%s\n%d 应该大于 %d", message, a, b)
	}
}

// AssertLessThan 断言a < b
func AssertLessThan(t *testing.T, a, b int, message string) {
	t.Helper()
	if a >= b {
		t.Errorf("%s\n%d 应该小于 %d", message, a, b)
	}
}
