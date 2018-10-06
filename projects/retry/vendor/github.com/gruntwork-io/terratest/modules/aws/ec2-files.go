package aws

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
)

// FetchContentsOfFileFromInstance looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
func FetchContentsOfFileFromInstance(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) string {
	out, err := FetchContentsOfFileFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// FetchContentsOfFileFromInstanceE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
func FetchContentsOfFileFromInstanceE(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) (string, error) {
	publicIp, err := GetPublicIpOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		return "", err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIp,
	}

	return ssh.FetchContentsOfFileE(t, host, useSudo, filePath)
}

// FetchContentsOfFilesFromInstance looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
func FetchContentsOfFilesFromInstance(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) map[string]string {
	out, err := FetchContentsOfFilesFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// FetchContentsOfFilesFromInstanceE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
func FetchContentsOfFilesFromInstanceE(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) (map[string]string, error) {
	publicIp, err := GetPublicIpOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		return nil, err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIp,
	}

	return ssh.FetchContentsOfFilesE(t, host, useSudo, filePaths...)
}

// FetchContentsOfFileFromAsg looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
func FetchContentsOfFileFromAsg(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) map[string]string {
	out, err := FetchContentsOfFileFromAsgE(t, awsRegion, sshUserName, keyPair, asgName, useSudo, filePath)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// FetchContentsOfFileFromAsgE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
func FetchContentsOfFileFromAsgE(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) (map[string]string, error) {
	instanceIDs, err := GetInstanceIdsForAsgE(t, asgName, awsRegion)
	if err != nil {
		return nil, err
	}

	instanceIdToContents := map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFileFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
		if err != nil {
			return nil, err
		}
		instanceIdToContents[instanceID] = contents
	}

	return instanceIdToContents, err
}

// FetchContentsOfFilesFromAsg looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
func FetchContentsOfFilesFromAsg(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) map[string]map[string]string {
	out, err := FetchContentsOfFilesFromAsgE(t, awsRegion, sshUserName, keyPair, asgName, useSudo, filePaths...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// FetchContentsOfFilesFromAsgE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
func FetchContentsOfFilesFromAsgE(t *testing.T, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) (map[string]map[string]string, error) {
	instanceIDs, err := GetInstanceIdsForAsgE(t, asgName, awsRegion)
	if err != nil {
		return nil, err
	}

	instanceIdToFilePathToContents := map[string]map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFilesFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
		if err != nil {
			return nil, err
		}
		instanceIdToFilePathToContents[instanceID] = contents
	}

	return instanceIdToFilePathToContents, err
}
