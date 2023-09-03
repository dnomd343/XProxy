package assets

import (
	. "XProxy/next/logger"
	"bytes"
	"github.com/andybalholm/brotli"
	"github.com/go-http-utils/headers"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"io"
	"net/http"
	"net/url"
)

// broltiDecode handles brolti encoding in http responses.
func broltiDecode(stream io.Reader) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, brotli.NewReader(stream))
	if err != nil {
		Logger.Errorf("Failed to decode http responses with brolti encoding -> %v", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

// gzipDecode handles gzip encoding in http responses.
func gzipDecode(stream io.Reader) ([]byte, error) {
	reader, err := gzip.NewReader(stream)
	if err != nil {
		Logger.Errorf("Failed to decode http responses with gzip encoding -> %v", err)
		return nil, err
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, reader)
	if err != nil {
		Logger.Errorf("Failed to handle gzip reader -> %v", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

// deflateDecode handles deflate encoding in http responses.
func deflateDecode(stream io.Reader) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, flate.NewReader(stream))
	if err != nil {
		Logger.Errorf("Failed to decode http responses with deflate encoding -> %v", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

// nonDecode handles plain encoding in http responses.
func nonDecode(stream io.Reader) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, stream)
	if err != nil {
		Logger.Errorf("Failed to read http responses -> %v", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

// createClient build http client based on http or socks proxy url.
func createClient(remoteUrl string, proxyUrl string) (http.Client, error) {
	if proxyUrl == "" {
		Logger.Infof("Downloading `%s` without proxy", remoteUrl)
		return http.Client{}, nil
	}
	Logger.Infof("Downloading `%s` via `%s`", remoteUrl, proxyUrl)
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		Logger.Errorf("Invalid proxy url `%s` -> %v", proxyUrl, err)
		return http.Client{}, err
	}
	return http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}, nil
}

// Download obtains resource file from the remote server and supports proxy.
func Download(url string, proxy string) ([]byte, error) {
	client, err := createClient(url, proxy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Errorf("Failed to create http request -> %v", err)
		return nil, err
	}
	req.Header.Set(headers.AcceptEncoding, "gzip, deflate, br")
	resp, err := client.Do(req)
	if err != nil {
		Logger.Errorf("Failed to execute http request -> %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	Logger.Debugf("Remote data downloaded successfully")

	var content []byte
	switch resp.Header.Get(headers.ContentEncoding) {
	case "br":
		Logger.Debugf("Downloaded content using brolti encoding")
		content, err = broltiDecode(resp.Body)
	case "gzip":
		Logger.Debugf("Downloaded content using gzip encoding")
		content, err = gzipDecode(resp.Body)
	case "deflate":
		Logger.Debugf("Downloaded content using deflate encoding")
		content, err = deflateDecode(resp.Body)
	default:
		content, err = nonDecode(resp.Body)
	}
	if err != nil {
		return nil, err
	}
	Logger.Debugf("Download `%s` successfully -> %d bytes", url, len(content))
	return content, nil
}
