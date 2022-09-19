package crd_generator

const CustomQuerySchema = `    Id: ID
    queryServiceTable(
        startTime: String
        endTime: String
        SystemServices: Boolean
        ShowGateways: Boolean
        Groupby: String
        noMetrics: Boolean
    ): TimeSeriesData
    queryServiceVersionTable(
        startTime: String
        endTime: String
        SystemServices: Boolean
        ShowGateways: Boolean
        noMetrics: Boolean
    ): TimeSeriesData
    queryServiceTS(
        svcMetric: String
        startTime: String
        endTime: String
        timeInterval: String
    ): TimeSeriesData
    queryIncomingAPIs(
        startTime: String
        endTime: String
        destinationService: String
        destinationServiceVersion: String
        timeInterval: String
        timeZone: String
    ): TimeSeriesData
    queryOutgoingAPIs(
        startTime: String
        endTime: String
        timeInterval: String
        timeZone: String
    ): TimeSeriesData
    queryIncomingTCP(
        startTime: String
        endTime: String
        destinationService: String
        destinationServiceVersion: String
    ): TimeSeriesData
    queryOutgoingTCP(
        startTime: String
        endTime: String
        destinationService: String
        destinationServiceVersion: String
    ): TimeSeriesData
    queryServiceTopology(
        metricStringArray: String
        startTime: String
        endTime: String
    ): TimeSeriesData
`
