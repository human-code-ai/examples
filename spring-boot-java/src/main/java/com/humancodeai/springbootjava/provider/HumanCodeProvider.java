package com.humancodeai.springbootjava.provider;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.humancodeai.springbootjava.config.HumanCodeConfig;
import com.humancodeai.springbootjava.model.ClientResponse;
import com.humancodeai.springbootjava.model.GetSessionIdResult;
import com.humancodeai.springbootjava.model.VerifyResult;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.*;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.HashMap;
import java.util.Map;

@Component
public class HumanCodeProvider {

    private final RestTemplate restTemplate;
    private final HumanCodeConfig config;
    private final ObjectMapper objectMapper = new ObjectMapper();

    public HumanCodeProvider(RestTemplateBuilder restTemplateBuilder, HumanCodeConfig config) {
        this.restTemplate = restTemplateBuilder
                .defaultHeader(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON_VALUE)
                .rootUri(config.getBaseUrl())
                .build();
        this.config = config;
        objectMapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
    }

    private String genSign(String data) throws NoSuchAlgorithmException, InvalidKeyException {
        Mac sha256 = Mac.getInstance("HmacSHA256");
        SecretKeySpec secretKey = new SecretKeySpec(config.getAppKey().getBytes(StandardCharsets.UTF_8), "HmacSHA256");
        sha256.init(secretKey);
        byte[] hashBytes = sha256.doFinal(data.getBytes(StandardCharsets.UTF_8));
        return bytesToHex(hashBytes);
    }

    private static String bytesToHex(byte[] bytes) {
        StringBuilder hexString = new StringBuilder();
        for (byte b : bytes) {
            String hex = String.format("%02x", b);
            hexString.append(hex);
        }
        return hexString.toString();
    }

    public GetSessionIdResult getSessionId(String nonceStr) throws Exception {
        long timestamp = System.currentTimeMillis();
        Map<String, String> requestBody = new HashMap<>();
        requestBody.put("timestamp", String.valueOf(timestamp));
        requestBody.put("nonce_str", nonceStr);

        String jsonBody = objectMapper.writeValueAsString(requestBody);
        String sign = genSign(jsonBody);

        String url = String.format("/api/session/v2/get_id?app_id=%s&sign=%s",
                config.getAppId(), sign);

        HttpEntity<String> entity = new HttpEntity<>(jsonBody);
        ResponseEntity<ClientResponse<GetSessionIdResult>> response = restTemplate.exchange(
                url, HttpMethod.POST, entity,
                new ParameterizedTypeReference<>() {
                });

        if (!response.getStatusCode().is2xxSuccessful() || response.getBody().getCode() != 0) {
            throw new RuntimeException("API error: " + response.getBody().getMsg());
        }
        return response.getBody().getResult();
    }

    public String genRegistrationUrl(String sessionId, String callBackUrl) {
        long timestamp = System.currentTimeMillis();
        return String.format("%s/authentication/index.html?session_id=%s&callback_url=%s&ts=%d#/",
                config.getBaseUrl(), sessionId, callBackUrl, timestamp);
    }

    public String genVerificationUrl(String sessionId, String humanId, String callBackUrl) {
        long timestamp = System.currentTimeMillis();
        return String.format("%s/authentication/index.html?session_id=%s&human_id=%s&callback_url=%s&ts=%d#/",
                config.getBaseUrl(), sessionId, humanId, callBackUrl, timestamp);
    }

    public VerifyResult verify(String sessionId, String vCode, String nonceStr) throws Exception {
        long timestamp = System.currentTimeMillis();
        Map<String, String> requestBody = new HashMap<>();
        requestBody.put("session_id", sessionId);
        requestBody.put("vcode", vCode);
        requestBody.put("timestamp", String.valueOf(timestamp));
        requestBody.put("nonce_str", nonceStr);

        String jsonBody = objectMapper.writeValueAsString(requestBody);
        String sign = genSign(jsonBody);

        String url = String.format("/api/vcode/v2/verify?app_id=%s&sign=%s",
                config.getAppId(), sign);

        HttpEntity<String> entity = new HttpEntity<>(jsonBody);
        ResponseEntity<ClientResponse<VerifyResult>> response = restTemplate.exchange(
                url, HttpMethod.POST, entity,
                new ParameterizedTypeReference<>() {
                });

        if (!response.getStatusCode().is2xxSuccessful() || response.getBody().getCode() != 0) {
            throw new RuntimeException("API error: " + response.getBody().getMsg());
        }
        return response.getBody().getResult();
    }
}
