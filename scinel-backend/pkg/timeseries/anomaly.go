package timeseries

import (
    "fmt"
    "math"
    "sort"
)

// AnomalyDetectionResult contains anomaly detection output
type AnomalyDetectionResult struct {
    IsAnomaly      bool
    ZScore         float64 // Number of standard deviations from mean
    Severity       string  // low, medium, high, critical
    HistoricalMean float64
    HistoricalStd  float64
    Details        string
}

// DetectAnomalyZScore detects anomalies using Z-score method
// Input: Historical values (time series), current value
// Output: Z-score and severity classification
// Reference: TDD Section 6.4.1 (Z-score method)
func DetectAnomalyZScore(historical []float64, currentValue float64) *AnomalyDetectionResult {
    if len(historical) < 5 {
        return &AnomalyDetectionResult{
            IsAnomaly: false,
            Details:   "Insufficient historical data",
        }
    }

    // Calculate mean
    mean := 0.0
    for _, v := range historical {
        mean += v
    }
    mean /= float64(len(historical))

    // Calculate standard deviation
    variance := 0.0
    for _, v := range historical {
        variance += (v - mean) * (v - mean)
    }
    stdDev := math.Sqrt(variance / float64(len(historical)-1))

    // Handle zero variance (all values identical)
    if stdDev == 0 {
        stdDev = math.Abs(mean) * 0.01 // Use 1% of mean as baseline
        if stdDev == 0 {
            stdDev = 1.0 // Fallback
        }
    }

    // Calculate Z-score
    zScore := (currentValue - mean) / stdDev

    // Determine anomaly and severity
    isAnomaly := math.Abs(zScore) > 2.0 // Z > 2.0 is significant
    severity := "low"
    if math.Abs(zScore) > 5.0 {
        severity = "critical"
    } else if math.Abs(zScore) > 4.0 {
        severity = "high"
    } else if math.Abs(zScore) > 3.0 {
        severity = "medium"
    }

    return &AnomalyDetectionResult{
        IsAnomaly:      isAnomaly,
        ZScore:         zScore,
        Severity:       severity,
        HistoricalMean: mean,
        HistoricalStd:  stdDev,
        Details:        fmt.Sprintf("Value %.1f is %.2f std deviations from mean %.1f", currentValue, zScore, mean),
    }
}

// DetectAnomalyIQR detects anomalies using Interquartile Range method
// Alternative method for outlier detection
// Reference: TDD Section 6.4.2 (IQR method)
func DetectAnomalyIQR(historical []float64, currentValue float64) *AnomalyDetectionResult {
    if len(historical) < 5 {
        return &AnomalyDetectionResult{
            IsAnomaly: false,
            Details:   "Insufficient historical data",
        }
    }

    // Sort values for quartile calculation
    sorted := make([]float64, len(historical))
    copy(sorted, historical)
    sort.Float64s(sorted)

    // Calculate Q1, Q2 (median), Q3
    q1Idx := len(sorted) / 4
    q3Idx := (3 * len(sorted)) / 4

    q1 := sorted[q1Idx]
    q3 := sorted[q3Idx]
    iqr := q3 - q1

    // Define outlier bounds
    lowerBound := q1 - 1.5*iqr
    upperBound := q3 + 1.5*iqr

    isAnomaly := currentValue < lowerBound || currentValue > upperBound

    // Determine severity based on distance from bounds
    var distance float64
    if currentValue < lowerBound {
        distance = lowerBound - currentValue
    } else if currentValue > upperBound {
        distance = currentValue - upperBound
    }

    severity := "low"
    if distance > iqr*3 {
        severity = "critical"
    } else if distance > iqr*2 {
        severity = "high"
    } else if distance > iqr {
        severity = "medium"
    }

    return &AnomalyDetectionResult{
        IsAnomaly:      isAnomaly,
        ZScore:         (currentValue - q1) / math.Max(iqr, 1.0), // Pseudo Z-score
        Severity:       severity,
        HistoricalMean: q1 + iqr/2, // Approximate mean
        HistoricalStd:  iqr / 1.35, // Approximate std dev from IQR
        Details:        fmt.Sprintf("Value %.1f outside IQR bounds [%.1f, %.1f]", currentValue, lowerBound, upperBound),
    }
}