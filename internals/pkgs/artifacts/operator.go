package artifacts

import (
	"context"
)

func (n *NabhikArtifactAgent) AppendFiles(ctx context.Context, files []string, source, destination *ArtifactOpts) {
	n.logger.Info().Msg("Appending files")
	// sourceCl := &auth.Client{
	// 	Header:     nil,
	// 	Credential: source.GetOrasCred(ctx),
	// }
	// destinationCl := &auth.Client{
	// 	Header:     nil,
	// 	Credential: destination.GetOrasCred(ctx),
	// }
	// sourceRepo, err := remote.NewRepository(source.ArtifactURL)
	// if err != nil {
	// 	n.logger.Error().Err(err).Msg("Error creating source artifact store from ORAS")
	// 	n.Error = err
	// 	return
	// }
	// destinationRepo, err := remote.NewRepository(destination.ArtifactURL)
	// if err != nil {
	// 	n.logger.Error().Err(err).Msg("Error creating destination artifact store from ORAS")
	// 	n.Error = err
	// 	return
	// }

	// sourceRepo.Client = sourceCl
	// destinationRepo.Client = destinationCl
	// memStore := memory.New()

	// sourceDesc, err := oras.Copy(ctx, sourceRepo, source.ArtifactURL, memStore, "", oras.DefaultCopyOptions)
	// if err != nil {
	// 	n.logger.Error().Err(err).Msg("Error copying source artifact to local memory store")
	// 	n.Error = err
	// 	return
	// }
	// n.logger.Info().Msgf("Source artifact copied to local memory store - %v", string(sourceDesc.Digest))
	// layers := []ocispec.Descriptor{}
	// for _, file := range files {
	// 	// Open the file
	// 	f, err := os.Open(file)
	// 	if err != nil {
	// 		fmt.Printf("Failed to open file %s: %v\n", file, err)
	// 		return
	// 	}
	// 	defer f.Close()

	// 	// Add the file to the store with the target path in the worker
	// 	destFile := filepath.Join(n.mountPath, filepath.Base(file))
	// 	descriptor := ocispec.Descriptor{
	// 		MediaType: ocispec.MediaTypeImageLayer,
	// 		ArtifactType: "layers",
	// 	}
	// 	err = memStore.Push(ctx, descriptor, f)
	// 	if err != nil {
	// 		fmt.Printf("Failed to add layer to store: %v\n", err)
	// 		return
	// 	}
	// 	layers = append(layers, layerDesc)
	// }

	// // Update the manifest with the new layers
	// manifestDesc, err := oras. (ctx, destinationRepo, destination.ArtifactURL, memStore, sourceDesc, layers, oras.DefaultAppendOptions)
	// if err != nil {
	// 	n.logger.Error().Err(err).Msgf("Error appending layers to the artifactusing ORAS for artifactURL - %v", source.ArtifactURL)
	// 	return
	// }

	// // Push the new image
	// destDesc, err := oras.Copy(ctx, memStore, "", destinationRepo, manifestDesc, oras.DefaultCopyOptions)
	// if err != nil {
	// 	n.logger.Error().Err(err).Msgf("Error copying artifact to destination using ORAS for artifactURL - %v", destination.ArtifactURL)
	// 	return
	// }
	// n.logger.Info().Msgf("Artifact copied to destination using ORAS for artifactURL - %v with digest - %v", destination.ArtifactURL, destDesc.Digest)

}

// func (n *NabhikArtifactAgent) AppendFiles(ctx context.Context, files []string, source, destination *ArtifactOpts) {
// 	var err error
// 	n.logger.Info().Msg("Appending files")
// 	sourceCl := &auth.Client{
// 		Header:     nil,
// 		Credential: source.GetOrasCred(ctx),
// 	}
// 	destinationCl := &auth.Client{
// 		Header:     nil,
// 		Credential: destination.GetOrasCred(ctx),
// 	}
// 	sourceRepo, err := remote.NewRepository(source.ArtifactURL)
// 	if err != nil {
// 		n.logger.Error().Err(err).Msg("Error creating source artifact store from ORAS")
// 		n.Error = err
// 		return
// 	}
// 	destinationRepo, err := remote.NewRepository(destination.ArtifactURL)
// 	if err != nil {
// 		n.logger.Error().Err(err).Msg("Error creating destination artifact store from ORAS")
// 		n.Error = err
// 		return
// 	}

// 	sourceRepo.Client = sourceCl
// 	destinationRepo.Client = destinationCl
// 	memStore := memory.New()

// 	sourceDesc, err := oras.Copy(ctx, sourceRepo, source.ArtifactURL, memStore, "", oras.DefaultCopyOptions)
// 	if err != nil {
// 		n.logger.Error().Err(err).Msg("Error copying source artifact to local memory store")
// 		n.Error = err
// 		return
// 	}
// 	n.logger.Info().Msgf("Source artifact copied to local memory store - %v", string(sourceDesc.Digest))
// 	layers := []ocispec.Descriptor{}
// 	for _, file := range files {
// 		// Open the file
// 		f, err := os.Open(file)
// 		if err != nil {
// 			fmt.Printf("Failed to open file %s: %v\n", file, err)
// 			return
// 		}
// 		defer f.Close()

// 		// Add the file to the store with the target path in the worker
// 		destFile := filepath.Join(n.mountPath, filepath.Base(file))
// 		descriptor := ocispec.Descriptor{
// 			MediaType: ocispec.MediaTypeImageLayer,
// 			ArtifactType: "layers",
// 		}
// 		err = memStore.Push(ctx, descriptor, f)
// 		if err != nil {
// 			fmt.Printf("Failed to add layer to store: %v\n", err)
// 			return
// 		}
// 		layers = append(layers, layerDesc)
// 	}

// 	// Update the manifest with the new layers
// 	manifestDesc, err := oras. (ctx, destinationRepo, destination.ArtifactURL, memStore, sourceDesc, layers, oras.DefaultAppendOptions)
// 	if err != nil {
// 		n.logger.Error().Err(err).Msgf("Error appending layers to the artifactusing ORAS for artifactURL - %v", source.ArtifactURL)
// 		return
// 	}

// 	// Push the new image
// 	destDesc, err := oras.Copy(ctx, memStore, "", destinationRepo, manifestDesc, oras.DefaultCopyOptions)
// 	if err != nil {
// 		n.logger.Error().Err(err).Msgf("Error copying artifact to destination using ORAS for artifactURL - %v", destination.ArtifactURL)
// 		return
// 	}
// 	n.logger.Info().Msgf("Artifact copied to destination using ORAS for artifactURL - %v with digest - %v", destination.ArtifactURL, destDesc.Digest)

// }
