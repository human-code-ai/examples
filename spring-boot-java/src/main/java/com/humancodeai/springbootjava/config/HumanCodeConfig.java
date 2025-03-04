package com.humancodeai.springbootjava.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Data
@ConfigurationProperties(prefix = "humancode")
@Component
public class HumanCodeConfig {
    String appId;
    String appKey;
    boolean debug;
    String baseUrl;
}
