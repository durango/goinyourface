test:
	go run image.go ./nfl.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m
	go run image.go ./face1.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m
	go run image.go ./face2.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m
	go run image.go ./face3.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m
	go run image.go ./face4.jpeg ./ccv/samples/pedestrian.icf ./ccv/samples/face ./ccv/samples/pedestrian.m
