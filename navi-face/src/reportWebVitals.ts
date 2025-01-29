import { onCLS, onFCP, onLCP, onTTFB } from 'web-vitals';

// Define the function to report web vitals
const reportWebVitals = (onPerfEntry?: (metric: any) => void) => {
  if (onPerfEntry && onPerfEntry instanceof Function) {
    // Pass the appropriate metric type to each handler
    onCLS(onPerfEntry);  // CLS
    onFCP(onPerfEntry);  // FCP
    onLCP(onPerfEntry);  // LCP
    onTTFB(onPerfEntry); // TTFB
  }
};

export default reportWebVitals;
